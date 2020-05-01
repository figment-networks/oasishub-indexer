package startpipeline

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"math/big"
)

type CalculatedValidatorData struct {
	EntityUID          types.PublicKey
	NodeUID            types.PublicKey
	ConsensusUID       types.PublicKey
	Address            string
	Proposed           bool
	VotingPower        types.VotingPower
	TotalShares        types.Quantity
	PrecommitValidated int64
	PrecommitType      *string
	PrecommitIndex     *int64
}

func CalculateValidatorsData(
	validatorsSyncable *syncable.Model,
	blockSyncable *syncable.Model,
	stateSyncable *syncable.Model,
) ([]CalculatedValidatorData, errors.ApplicationError) {
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

	var validators []CalculatedValidatorData
	for i, rv := range validatorsData {
		e := CalculatedValidatorData{
			EntityUID:    types.PublicKey(rv.Node.EntityId),
			NodeUID:      types.PublicKey(rv.Id),
			ConsensusUID: types.PublicKey(rv.Node.Consensus.Id),
			Address:      rv.Address,
			VotingPower:  types.VotingPower(rv.VotingPower),
		}

		// Get precommit data
		var validated int64
		var index int64
		var pType string
		if len(blockData.LastCommit.Votes) > 0 {
			// Account for situation when there is more validators than precommits
			// It means that last x validators did not have chance to vote. In that case set validated to null.
			if i > len(blockData.LastCommit.Votes)-1 {
				index = int64(i)
				validated = 2
			} else {
				precommit := blockData.LastCommit.Votes[i]

				if precommit == nil {
					validated = 0
					index = int64(i)
				} else {
					validated = 1
					index = precommit.ValidatorIndex
					pType = precommit.Type
				}
			}
		} else {
			validated = 2
		}

		e.PrecommitValidated = validated
		e.PrecommitIndex = &index
		e.PrecommitType = &pType

		// Get proposed
		e.Proposed = blockData.Header.ProposerAddress == e.Address

		// Get total shares
		delegations := stateData.Staking.Delegations[rv.Node.EntityId]
		totalShares := big.NewInt(0)
		for _, d := range delegations.Entries {
			shares := types.NewQuantityFromBytes(d.Shares)
			totalShares = totalShares.Add(totalShares, &shares.Int)
		}
		e.TotalShares = types.NewQuantity(totalShares)

		validators = append(validators, e)
	}
	return validators, nil
}

type CalculatedAccountData struct {
	PublicKey                         types.PublicKey
	RecentGeneralBalance             types.Quantity
	RecentGeneralNonce               types.Nonce
	RecentEscrowActiveBalance        types.Quantity
	RecentEscrowActiveTotalShares    types.Quantity
	RecentEscrowDebondingBalance     types.Quantity
	RecentEscrowDebondingTotalShares types.Quantity
}

func CalculateAccountsData(stateSyncable *syncable.Model) ([]CalculatedAccountData, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	var accounts []CalculatedAccountData
	for publicKey, info := range stateData.Staking.Ledger {
		acc := CalculatedAccountData{
			PublicKey:                         types.PublicKey(publicKey),
			RecentGeneralBalance:             types.NewQuantityFromBytes(info.General.Balance),
			RecentGeneralNonce:               types.Nonce(info.General.Nonce),
			RecentEscrowActiveBalance:        types.NewQuantityFromBytes(info.Escrow.Active.Balance),
			RecentEscrowActiveTotalShares:    types.NewQuantityFromBytes(info.Escrow.Active.TotalShares),
			RecentEscrowDebondingBalance:     types.NewQuantityFromBytes(info.Escrow.Debonding.Balance),
			RecentEscrowDebondingTotalShares: types.NewQuantityFromBytes(info.Escrow.Debonding.TotalShares),
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}
