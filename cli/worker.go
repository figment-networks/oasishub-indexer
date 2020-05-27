package cli

import (
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/usecase"
	"github.com/figment-networks/oasishub-indexer/worker"
)

func startWorker(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()
	client, err := initClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	workerHandlers := usecase.NewWorkerHandlers(cfg, db, client)

	w, err := worker.New(cfg, workerHandlers)
	if err != nil {
		return err
	}

	w.Start()

	return nil
}
