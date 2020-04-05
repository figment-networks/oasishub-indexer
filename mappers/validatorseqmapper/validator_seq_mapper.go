package validatorseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
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
			if i > len(blockData.Data.LastCommit.Precommits) - 1 {
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

func ToView(ts []*validatorseq.Model) []map[string]interface{} {
	var items []map[string]interface{}
	for _, t := range ts {
		i := map[string]interface{}{
			"id":       t.ID,
			"height":   t.Height,
			"time":     t.Time,
			"chain_id": t.ChainId,

			"entity_uid":    t.EntityUID,
			"node_uid":      t.NodeUID,
			"consensus_uid": t.ConsensusUID,
			"voting_power":  t.VotingPower,
			"proposed":      t.Proposed,
		}
		items = append(items, i)
	}
	return items
}
