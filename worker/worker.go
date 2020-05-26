package worker

import (
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/usecase"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/robfig/cron/v3"
)

var (
	job cron.Job
)

type Worker struct {
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

	job = cron.FuncJob(handlers.RunIndexer.Handle)
	job = cron.NewChain(cron.SkipIfStillRunning(log)).Then(job)
	_, err := cronJob.AddJob(cfg.SyncInterval, job)
	if err != nil {
		return nil, err
	}

	return &Worker{
		cronJob: cronJob,
	}, nil
}

func (w *Worker) Start() {
	logger.Info("starting worker...", logger.Field("app", "worker"))

	w.cronJob.Start()

	//Run forever
	select {}
}
