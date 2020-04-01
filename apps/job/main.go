package main

import (
	"github.com/robfig/cron/v3"
)

var (
	cronJob *cron.Cron
)

func main() {
	// CLIENTS
	//node := shared.NewNodeClient()
	//db := shared.NewDbClient()

	// REPOSITORIES
	//blockDbRepo := blockseqrepo.NewDbRepo(db.Client())
	//reportDbRepo := reportrepo.NewDbRepo(db.Client())

	//USE CASES

	// HANDLERS

	// ADD CRON FUNCS
	//_, err := cronJob.AddFunc(config.SyncInterval(), startBlockPipelineHandler.Handle)
	//if err != nil {
	//	log.Error(err)
	//	panic(err)
	//}

	// START CRON
	cronJob.Start()
	//Run forever
	select {}
}