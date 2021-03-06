package worker

import (
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/usecase"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/figment-networks/oasishub-indexer/utils/reporting"
	"github.com/robfig/cron/v3"
)

var (
	job cron.Job
)

type Worker struct {
	cfg      *config.Config
	handlers *usecase.WorkerHandlers

	logger  logger.CronLogger
	cronJob *cron.Cron
}

func New(cfg *config.Config, handlers *usecase.WorkerHandlers) (*Worker, error) {
	log := logger.NewCronLogger()
	cronJob := cron.New(
		cron.WithLogger(cron.VerbosePrintfLogger(log)),
		cron.WithChain(
			cron.Recover(log),
		),
	)

	w := &Worker{
		cfg:      cfg,
		handlers: handlers,
		logger:   log,
		cronJob:  cronJob,
	}

	return w.init()
}

func (w *Worker) init() (*Worker, error) {
	_, err := w.addIndexerIndexJob()
	if err != nil {
		return nil, err
	}

	_, err = w.addIndexerSummarizeJob()
	if err != nil {
		return nil, err
	}

	_, err = w.addIndexerPurgeJob()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Worker) Start() error {
	defer reporting.RecoverError()

	logger.Info("starting worker...", logger.Field("app", "worker"))

	w.cronJob.Start()

	return nil
}
