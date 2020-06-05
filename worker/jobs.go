package worker

import "github.com/robfig/cron/v3"

func (w *Worker) addRunIndexerJob() (cron.EntryID, error) {
	job = cron.FuncJob(w.handlers.RunIndexer.Handle)
	job = cron.NewChain(cron.SkipIfStillRunning(w.logger)).Then(job)
	return w.cronJob.AddJob(w.cfg.SyncInterval, job)
}

func (w *Worker) addPurgeIndexerJob() (cron.EntryID, error) {
	job = cron.FuncJob(w.handlers.PurgeIndexer.Handle)
	job = cron.NewChain(cron.SkipIfStillRunning(w.logger)).Then(job)
	return w.cronJob.AddJob(w.cfg.PurgeInterval, job)
}
