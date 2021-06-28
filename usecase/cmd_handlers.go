package usecase

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/usecase/chain"
	"github.com/figment-networks/oasishub-indexer/usecase/indexing"
	"github.com/figment-networks/oasishub-indexer/usecase/validator"
)

func NewCmdHandlers(cfg *config.Config, db *store.Store, c *client.Client) *CmdHandlers {
	return &CmdHandlers{
		GetStatus:          chain.NewGetStatusCmdHandler(db, c),
		IndexerIndex:       indexing.NewIndexCmdHandler(cfg, db, c),
		IndexerBackfill:    indexing.NewBackfillCmdHandler(cfg, db, c),
		IndexerPurge:       indexing.NewPurgeCmdHandler(cfg, db, c),
		IndexerReindex:     indexing.NewReindexCmdHandler(cfg, db, c),
		IndexerSummarize:   indexing.NewSummarizeCmdHandler(cfg, db, c),
		DecorateValidators: validator.NewDecorateCmdHandler(cfg, db, c),
	}
}

type CmdHandlers struct {
	GetStatus          *chain.GetStatusCmdHandler
	IndexerIndex       *indexing.IndexCmdHandler
	IndexerBackfill    *indexing.BackfillCmdHandler
	IndexerPurge       *indexing.PurgeCmdHandler
	IndexerReindex     *indexing.ReindexCmdHandler
	IndexerSummarize   *indexing.SummarizeCmdHandler
	DecorateValidators *validator.DecorateCmdHandler
}
