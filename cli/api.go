package cli

import (
	"github.com/figment-networks/oasishub-indexer/api"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/usecase"
)

func startApi(cfg *config.Config) error {
	client, err := initClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	httpHandlers := usecase.NewHttpHandlers(db, client)

	a := api.New(httpHandlers)
	if err := a.Start(cfg.ListenAddr()); err != nil {
		return err
	}
	return nil
}
