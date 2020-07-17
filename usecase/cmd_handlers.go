package usecase

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/indexing"
)

func NewCmdHandlers(cfg *config.Config, db *store.Store, c *client.Client) *CmdHandlers {
	return &CmdHandlers{
		RunIndexer:       indexing.NewRunCmdHandler(cfg, db, c),
		PurgeIndexer:     indexing.NewPurgeCmdHandler(cfg, db, c),
		SummarizeIndexer: indexing.NewSummarizeCmdHandler(cfg, db, c),
		ParseCSV:         indexing.NewParseCSVCmdHandler(cfg, db, c),
	}
}

type CmdHandlers struct {
	RunIndexer       types.CmdHandler
	PurgeIndexer     types.CmdHandler
	SummarizeIndexer types.CmdHandler
	ParseCSV         types.CmdHandler
}
