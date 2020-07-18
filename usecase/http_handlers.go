package usecase

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/account"
	"github.com/figment-networks/oasishub-indexer/usecase/block"
	"github.com/figment-networks/oasishub-indexer/usecase/chain"
	"github.com/figment-networks/oasishub-indexer/usecase/debondingdelegation"
	"github.com/figment-networks/oasishub-indexer/usecase/delegation"
	"github.com/figment-networks/oasishub-indexer/usecase/health"
	"github.com/figment-networks/oasishub-indexer/usecase/staking"
	"github.com/figment-networks/oasishub-indexer/usecase/systemevent"
	"github.com/figment-networks/oasishub-indexer/usecase/transaction"
	"github.com/figment-networks/oasishub-indexer/usecase/validator"
)

func NewHttpHandlers(cfg *config.Config, db *store.Store, c *client.Client) *HttpHandlers {
	return &HttpHandlers{
		Health:                          health.NewHealthHttpHandler(),
		GetStatus:                       chain.NewGetStatusHttpHandler(db, c),
		GetBlockByHeight:                block.NewGetByHeightHttpHandler(db, c),
		GetBlockTimes:                   block.NewGetBlockTimesHttpHandler(db, c),
		GetBlockSummary:                 block.NewGetBlockSummaryHttpHandler(db, c),
		GetAccountByAddress:             account.NewGetByAddressHttpHandler(db, c),
		GetDebondingDelegationsByHeight: debondingdelegation.NewGetByHeightHttpHandler(db, c),
		GetDelegationsByHeight:          delegation.NewGetByHeightHttpHandler(db, c),
		GetStakingDetailsByHeight:       staking.NewGetByHeightHttpHandler(db, c),
		GetTransactionsByHeight:         transaction.NewGetByHeightHttpHandler(db, c),
		GetValidatorsByHeight:           validator.NewGetByHeightHttpHandler(cfg, db, c),
		GetValidatorByAddress:           validator.NewGetByAddressHttpHandler(db, c),
		GetValidatorSummary:             validator.NewGetSummaryHttpHandler(db, c),
		GetValidatorsForMinHeight:       validator.NewGetForMinHeightHttpHandler(db, c),
		GetSystemEventsForAddress:       systemevent.NewGetForAddressHttpHandler(db, c),
	}
}

type HttpHandlers struct {
	Health                          types.HttpHandler
	GetStatus                       types.HttpHandler
	GetBlockTimes                   types.HttpHandler
	GetBlockSummary                 types.HttpHandler
	GetBlockByHeight                types.HttpHandler
	GetAccountByAddress             types.HttpHandler
	GetDebondingDelegationsByHeight types.HttpHandler
	GetDelegationsByHeight          types.HttpHandler
	GetStakingDetailsByHeight       types.HttpHandler
	GetTransactionsByHeight         types.HttpHandler
	GetValidatorsByHeight           types.HttpHandler
	GetValidatorByAddress           types.HttpHandler
	GetValidatorSummary             types.HttpHandler
	GetValidatorsForMinHeight       types.HttpHandler
	GetSystemEventsForAddress       types.HttpHandler
}
