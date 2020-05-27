package usecase

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
)

func New(db *store.Store, c *client.Client) *UseCase {
	return &UseCase{
		HttpHandlers: NewHttpHandlers(db, c),
	}
}

type UseCase struct {
	HttpHandlers   *HttpHandlers
	CLIHandlers    *CmdHandlers
	WorkerHandlers *WorkerHandlers
}