package usecase

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/account"
	"github.com/figment-networks/oasishub-indexer/usecase/block"
	"github.com/figment-networks/oasishub-indexer/usecase/debondingdelegation"
	"github.com/figment-networks/oasishub-indexer/usecase/delegation"
	"github.com/figment-networks/oasishub-indexer/usecase/health"
	"github.com/figment-networks/oasishub-indexer/usecase/staking"
	"github.com/figment-networks/oasishub-indexer/usecase/syncable"
	"github.com/figment-networks/oasishub-indexer/usecase/transaction"
	"github.com/figment-networks/oasishub-indexer/usecase/validator"
)

func NewHttpHandlers(db *store.Store, c *client.Client) *HttpHandlers {
	return &HttpHandlers{
		Health:                          health.NewHealthHttpHandler(),
		GetBlockByHeight:                block.NewGetByHeightHttpHandler(db, c),
		GetBlockTimes:                   block.NewGetBlockTimesHttpHandler(db, c),
		GetBlockSummary:                 block.NewGetBlockSummaryHttpHandler(db, c),
		GetAccountByPublicKey:           account.NewGetByPublicKeyHttpHandler(db, c),
		GetDebondingDelegationsByHeight: debondingdelegation.NewGetByHeightHttpHandler(db, c),
		GetDelegationsByHeight:          delegation.NewGetByHeightHttpHandler(db, c),
		GetStakingDetailsByHeight:       staking.NewGetByHeightHttpHandler(db, c),
		GetMostRecentHeight:             syncable.NewGetMostRecentHeightHttpHandler(db, c),
		GetTransactionsByHeight:         transaction.NewGetByHeightHttpHandler(db, c),
		GetValidatorsByHeight:           validator.NewGetByHeightHttpHandler(db, c),
		GetValidatorByEntityUid:         validator.NewGetByEntityUidHttpHandler(db, c),
		GetValidatorSummary:             validator.NewGetSummaryHttpHandler(db, c),
		GetValidatorsForMinHeight:       validator.NewGetForMinHeightHttpHandler(db, c),
	}
}

type HttpHandlers struct {
	Health                          types.HttpHandler
	GetBlockTimes                   types.HttpHandler
	GetBlockSummary                 types.HttpHandler
	GetBlockByHeight                types.HttpHandler
	GetAccountByPublicKey           types.HttpHandler
	GetDebondingDelegationsByHeight types.HttpHandler
	GetDelegationsByHeight          types.HttpHandler
	GetStakingDetailsByHeight       types.HttpHandler
	GetMostRecentHeight             types.HttpHandler
	GetTransactionsByHeight         types.HttpHandler
	GetValidatorsByHeight           types.HttpHandler
	GetValidatorByEntityUid         types.HttpHandler
	GetValidatorSummary             types.HttpHandler
	GetValidatorsForMinHeight       types.HttpHandler
}
