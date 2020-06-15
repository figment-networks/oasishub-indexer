package usecase

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/usecase/indexing"
)

func NewCmdHandlers(cfg *config.Config, db *store.Store, c *client.Client) *CmdHandlers {
	return &CmdHandlers{
		StartIndexer:     indexing.NewStartCmdHandler(cfg, db, c),
		BackfillIndexer:  indexing.NewBackfillCmdHandler(cfg, db, c),
		PurgeIndexer:     indexing.NewPurgeCmdHandler(cfg, db, c),
		SummarizeIndexer: indexing.NewSummarizeCmdHandler(cfg, db, c),
	}
}

type CmdHandlers struct {
	StartIndexer     *indexing.StartCmdHandler
	BackfillIndexer  *indexing.BackfillCmdHandler
	PurgeIndexer     *indexing.PurgeCmdHandler
	SummarizeIndexer *indexing.SummarizeCmdHandler
}
