package worker

import "github.com/robfig/cron/v3"

func (w *Worker) addIndexerIndexJob() (cron.EntryID, error) {
	job = cron.FuncJob(w.handlers.IndexerIndex.Handle)
	job = cron.NewChain(cron.SkipIfStillRunning(w.logger)).Then(job)
	return w.cronJob.AddJob(w.cfg.IndexWorkerInterval, job)
}

func (w *Worker) addIndexerSummarizeJob() (cron.EntryID, error) {
	job = cron.FuncJob(w.handlers.IndexerSummarize.Handle)
	job = cron.NewChain(cron.SkipIfStillRunning(w.logger)).Then(job)
	return w.cronJob.AddJob(w.cfg.SummarizeWorkerInterval, job)
}

func (w *Worker) addIndexerPurgeJob() (cron.EntryID, error) {
	job = cron.FuncJob(w.handlers.IndexerPurge.Handle)
	job = cron.NewChain(cron.SkipIfStillRunning(w.logger)).Then(job)
	return w.cronJob.AddJob(w.cfg.PurgeWorkerInterval, job)
}
