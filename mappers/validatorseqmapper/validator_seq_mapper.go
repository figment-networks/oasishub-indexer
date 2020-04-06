package validatorseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/debondingdelegationseq"
	"github.com/figment-networks/oasishub-indexer/models/delegationseq"
	"github.com/figment-networks/oasishub-indexer/models/entityagg"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/models/validatorseq"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"math/big"
)

func ToSequence(validatorsSyncable syncable.Model, blockSyncable syncable.Model, stateSyncable syncable.Model) ([]validatorseq.Model, errors.ApplicationError) {
	validatorsData, err := syncablemapper.UnmarshalValidatorsData(validatorsSyncable.Data)
	if err != nil {
		return nil, err
	}
	blockData, err := syncablemapper.UnmarshalBlockData(blockSyncable.Data)
	if err != nil {
		return nil, err
	}
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	var validators []validatorseq.Model
	for i, rv := range validatorsData.Data {
		e := validatorseq.Model{
			Sequence: &shared.Sequence{
				ChainId: validatorsSyncable.ChainId,
				Height:  validatorsSyncable.Height,
				Time:    validatorsSyncable.Time,
			},

			EntityUID:    types.PublicKey(rv.Node.EntityID.String()),
			NodeUID:      types.PublicKey(rv.ID.String()),
			ConsensusUID: types.PublicKey(rv.Node.Consensus.ID.String()),
			Address:      rv.Address,
			VotingPower:  validatorseq.VotingPower(rv.VotingPower),
		}

		// Get precommit data
		if len(blockData.Data.LastCommit.Precommits) > 0 {
			var validated bool
			var index int64
			var pType int64
			// Account for situation when there is more validators than precommits
			// It means that last x validators did not have chance to vote. In that case set validated to null.
			if i > len(blockData.Data.LastCommit.Precommits)-1 {
				index = int64(i)
			} else {
				precommit := blockData.Data.LastCommit.Precommits[i]

				if precommit == nil {
					validated = false
					index = int64(i)
				} else {
					validated = true
					index = int64(precommit.ValidatorIndex)
					pType = int64(precommit.Type)
				}
			}
			e.PrecommitValidated = &validated
			e.PrecommitIndex = &index
			e.PrecommitType = &pType
		}

		// Get proposed
		e.Proposed = blockData.Data.Header.ProposerAddress.String() == e.Address

		// Get total shares
		delegations := stateData.Data.Staking.Delegations[rv.Node.EntityID]
		totalShares := big.NewInt(0)
		for _, d := range delegations {
			totalShares = totalShares.Add(totalShares, d.Shares.ToBigInt())
		}
		e.TotalShares = types.NewQuantity(totalShares)

		if !e.Valid() {
			return nil, errors.NewErrorFromMessage("validator sequence not valid", errors.NotValid)
		}

		validators = append(validators, e)
	}
	return validators, nil
}

type ListView struct {
	Items []validatorseq.Model `json:"items"`
}

func ToListView(ms []validatorseq.Model) *ListView {
	return &ListView{
		Items: ms,
	}
}

type DetailsView struct {
	*shared.Model
	*shared.Aggregate

	EntityUID types.PublicKey `json:"entity_uid"`

	TotalValidated             int64                          `json:"total_validated"`
	TotalMissed                int64                          `json:"total_missed"`
	TotalProposed              int64                          `json:"total_proposed"`
	CurrentDelegations         []delegationseq.Model          `json:"current_delegations"`
	RecentDebondingDelegations []debondingdelegationseq.Model `json:"recent_debonding_delegations"`
}

func ToDetailsView(m entityagg.Model, totV int64, totM int64, totP int64, currDs []delegationseq.Model, recDds []debondingdelegationseq.Model) *DetailsView {
	return &DetailsView{
		Model:     m.Model,
		Aggregate: m.Aggregate,

		TotalValidated: totV,
		TotalMissed:    totM,
		TotalProposed:  totP,

		CurrentDelegations:         currDs,
		RecentDebondingDelegations: recDds,
	}
}
