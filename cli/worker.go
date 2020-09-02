package cli

import (
	"net/http"
	"time"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/metrics/prometheusmetrics"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/usecase"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
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

	prom := prometheusmetrics.New()
	err = metrics.AddEngine(prom)
	if err != nil {
		logger.Error(err)
	}
	err = metrics.Hotload(prom.Name())
	if err != nil {
		logger.Error(err)
	}
	s := &http.Server{
		Addr:         cfg.IndexerMetricAddr,
		Handler:      metrics.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s.ListenAndServe()
}
