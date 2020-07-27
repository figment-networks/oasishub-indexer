package chain

import (
	"context"
	"fmt"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

type GetStatusCmdHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getStatusUseCase
}

func NewGetStatusCmdHandler(db *store.Store, c *client.Client) *GetStatusCmdHandler {
	return &GetStatusCmdHandler{
		db:     db,
		client: c,
	}
}

func (h *GetStatusCmdHandler) Handle(ctx context.Context) {
	logger.Info("chain get status use case [handler=cmd]")

	details, err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}

	fmt.Println("=== App ===")
	fmt.Println("Name:", details.AppName)
	fmt.Println("Version", details.AppVersion)
	fmt.Println("Go version", details.GoVersion)
	fmt.Println("")

	fmt.Println("=== Chain ===")
	fmt.Println("ID:", details.ChainID)
	fmt.Println("Name:", details.ChainName)
	fmt.Println("App version:", details.ChainAppVersion)
	fmt.Println("Block version:", details.ChainBlockVersion)
	fmt.Println("")

	fmt.Println("=== Genesis ===")
	fmt.Println("Height:", details.GenesisHeight)
	fmt.Println("Time:", details.GenesisTime)
	fmt.Println("")

	fmt.Println("=== Indexing ===")
	fmt.Println("Last index version:", details.LastIndexVersion)
	fmt.Println("Last indexed height:", details.LastIndexedHeight)
	fmt.Println("Last indexed time:", details.LastIndexedTime)
	fmt.Println("Last indexed at:", details.LastIndexedAt)
	fmt.Println("Lag behind head:", details.Lag)
	fmt.Println("")
}

func (h *GetStatusCmdHandler) getUseCase() *getStatusUseCase {
	if h.useCase == nil {
		h.useCase = NewGetStatusUseCase(h.db, h.client)
	}
	return h.useCase
}

