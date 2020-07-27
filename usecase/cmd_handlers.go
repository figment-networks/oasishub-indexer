package usecase

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/usecase/chain"
	"github.com/figment-networks/oasishub-indexer/usecase/indexing"
)

func NewCmdHandlers(cfg *config.Config, db *store.Store, c *client.Client) *CmdHandlers {
	return &CmdHandlers{
		GetStatus:        chain.NewGetStatusCmdHandler(db, c),
		IndexerIndex:     indexing.NewIndexCmdHandler(cfg, db, c),
		IndexerBackfill:  indexing.NewBackfillCmdHandler(cfg, db, c),
		IndexerPurge:     indexing.NewPurgeCmdHandler(cfg, db, c),
		IndexerSummarize: indexing.NewSummarizeCmdHandler(cfg, db, c),
	}
}

type CmdHandlers struct {
	GetStatus        *chain.GetStatusCmdHandler
	IndexerIndex     *indexing.IndexCmdHandler
	IndexerBackfill  *indexing.BackfillCmdHandler
	IndexerPurge     *indexing.PurgeCmdHandler
	IndexerSummarize *indexing.SummarizeCmdHandler
}
