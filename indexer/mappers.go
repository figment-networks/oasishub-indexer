package indexer

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/pkg/errors"
)

var (
	errInvalidBlockSeq = errors.New("block sequence not valid")
)

func BlockToSequence(syncable *model.Syncable, blockParsedData ParsedBlockData) (*model.BlockSeq, error) {
	e := &model.BlockSeq{
		Sequence: &model.Sequence{
			Height: syncable.Height,
			Time:   syncable.Time,
		},

		TransactionsCount: blockParsedData.TransactionsCount,
	}

	if !e.Valid() {
		return nil, errInvalidBlockSeq
	}

	return e, nil
}

func ValidatorToSequence(syncable *model.Syncable, rawValidators []*validatorpb.Validator, parsedValidators ParsedValidatorsData) ([]model.ValidatorSeq, error) {
	var validators []model.ValidatorSeq
	for _, rawValidator := range rawValidators {
		e := model.ValidatorSeq{
			Sequence: &model.Sequence{
				Height: syncable.Height,
				Time:   syncable.Time,
			},

			EntityUID:   rawValidator.GetNode().GetEntityId(),
			Address:     rawValidator.GetAddress(),
			VotingPower: rawValidator.GetVotingPower(),
			Commission:  types.NewQuantityFromBytes(rawValidator.GetCommission()),
		}

		address := rawValidator.GetAddress()
		parsedValidator, ok := parsedValidators[address]
		if ok {
			e.PrecommitValidated = parsedValidator.PrecommitValidated
			e.Proposed = parsedValidator.Proposed
			e.TotalShares = parsedValidator.TotalShares
			e.ActiveEscrowBalance = parsedValidator.ActiveEscrowBalance
		}

		if !e.Valid() {
			return nil, errors.New("validator sequence not valid")
		}

		validators = append(validators, e)
	}
	return validators, nil
}

func TransactionToSequence(syncable *model.Syncable, rawTransactions []*transactionpb.Transaction) ([]model.TransactionSeq, error) {
	var transactions []model.TransactionSeq
	for _, rawTransaction := range rawTransactions {
		e := model.TransactionSeq{
			Sequence: &model.Sequence{
				Height: syncable.Height,
				Time:   syncable.Time,
			},

			PublicKey: rawTransaction.GetPublicKey(),
			Hash:      rawTransaction.GetHash(),
			Nonce:     rawTransaction.GetNonce(),
			Fee:       types.NewQuantityFromBytes(rawTransaction.GetFee()),
			GasLimit:  rawTransaction.GetGasLimit(),
			GasPrice:  types.NewQuantityFromBytes(rawTransaction.GetGasPrice()),
			Method:    rawTransaction.GetMethod(),
		}

		if !e.Valid() {
			return nil, errors.New("transaction sequence not valid")
		}

		transactions = append(transactions, e)
	}
	return transactions, nil
}

func StakingToSequence(syncable *model.Syncable, rawStaking *statepb.Staking) (*model.StakingSeq, error) {
	e := &model.StakingSeq{
		Sequence: &model.Sequence{
			Height: syncable.Height,
			Time:   syncable.Time,
		},

		TotalSupply:         types.NewQuantityFromBytes(rawStaking.GetTotalSupply()),
		CommonPool:          types.NewQuantityFromBytes(rawStaking.GetCommonPool()),
		DebondingInterval:   rawStaking.GetParameters().GetDebondingInterval(),
		MinDelegationAmount: types.NewQuantityFromBytes(rawStaking.GetParameters().GetMinDelegationAmount()),
	}

	if !e.Valid() {
		return nil, errors.New("staking sequence not valid")
	}

	return e, nil
}

func DelegationToSequence(syncable *model.Syncable, rawState *statepb.State) ([]model.DelegationSeq, error) {
	var delegations []model.DelegationSeq
	for validatorUID, delegationsMap := range rawState.GetStaking().GetDelegations() {
		for delegatorUID, info := range delegationsMap.GetEntries() {
			acc := model.DelegationSeq{
				Sequence: &model.Sequence{
					Height: syncable.Height,
					Time:   syncable.Time,
				},

				ValidatorUID: validatorUID,
				DelegatorUID: delegatorUID,
				Shares:       types.NewQuantityFromBytes(info.GetShares()),
			}

			if !acc.Valid() {
				return nil, errors.New("delegation sequence not valid")
			}

			delegations = append(delegations, acc)
		}
	}
	return delegations, nil
}

func DebondingDelegationToSequence(syncable *model.Syncable, rawState *statepb.State) ([]model.DebondingDelegationSeq, error) {
	var delegations []model.DebondingDelegationSeq
	for validatorUID, delegationsMap := range rawState.GetStaking().GetDebondingDelegations() {
		for delegatorUID, infoArray := range delegationsMap.GetEntries() {
			for _, delegation := range infoArray.GetDebondingDelegations() {
				acc := model.DebondingDelegationSeq{
					Sequence: &model.Sequence{
						Height: syncable.Height,
						Time:   syncable.Time,
					},

					ValidatorUID: validatorUID,
					DelegatorUID: delegatorUID,
					Shares:       types.NewQuantityFromBytes(delegation.GetShares()),
					DebondEnd:    delegation.GetDebondEndTime(),
				}

				if !acc.Valid() {
					return nil, errors.New("debonding delegation sequence not valid")
				}

				delegations = append(delegations, acc)
			}
		}
	}
	return delegations, nil
}
