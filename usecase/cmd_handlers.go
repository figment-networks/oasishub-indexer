package usecase

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/indexer"
)

func NewCmdHandlers(cfg *config.Config, db *store.Store, c *client.Client) *CmdHandlers {
	return &CmdHandlers{
		RunIndexer: indexer.NewRunIndexerCmdHandler(cfg, db, c),
	}
}

type CmdHandlers struct {
	RunIndexer types.CmdHandler
}
