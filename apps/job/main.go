package main

import (
	"github.com/figment-networks/oasishub-indexer/apps/shared"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/entityaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/reportrepo"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/usecases/syncable/cleanup"
	"github.com/figment-networks/oasishub-indexer/usecases/syncable/startpipeline"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/robfig/cron/v3"
)

var (
	cronJob *cron.Cron
	cronLog = log.NewCronLogger()
	job     cron.Job

)

func main() {
	// CLIENTS
	node := shared.NewNodeClient()
	db := shared.NewDbClient()

	// REPOSITORIES
	syncableDbRepo := syncablerepo.NewDbRepo(db.Client())
	syncableProxyRepo := syncablerepo.NewProxyRepo(node)
	reportDbRepo := reportrepo.NewDbRepo(db.Client())

	blockSeqDbRepo := blockseqrepo.NewDbRepo(db.Client())
	validatorSequenceDbRepo := validatorseqrepo.NewDbRepo(db.Client())
	transactionSeqDbRepo := transactionseqrepo.NewDbRepo(db.Client())
	stakingSeqDbRepo := stakingseqrepo.NewDbRepo(db.Client())
	delegationSeqDbRepo := delegationseqrepo.NewDbRepo(db.Client())
	debondingDelegationSeqDbRepo := debondingdelegationseqrepo.NewDbRepo(db.Client())

	accountAggDbRepo := accountaggrepo.NewDbRepo(db.Client())
	entityAggDbRepo := entityaggrepo.NewDbRepo(db.Client())

	//USE CASES
	startPipelineUseCase := startpipeline.NewUseCase(
		syncableDbRepo,
		syncableProxyRepo,
		blockSeqDbRepo,
		validatorSequenceDbRepo,
		transactionSeqDbRepo,
		stakingSeqDbRepo,
		accountAggDbRepo,
		delegationSeqDbRepo,
		debondingDelegationSeqDbRepo,
		entityAggDbRepo,
		reportDbRepo,
	)
	cleanupUseCase := cleanup.NewUseCase(syncableDbRepo)

	// HANDLERS
	startPipelineHandler := startpipeline.NewJobHandler(startPipelineUseCase)
	cleanupHandler := cleanup.NewJobHandler(cleanupUseCase)

	// CRON
	cronJob = cron.New(
		cron.WithLogger(cron.VerbosePrintfLogger(log.GetLogger())),
		cron.WithChain(
			cron.Recover(cronLog),
		),
	)

	// Start pipeline job
	job = cron.FuncJob(startPipelineHandler.Handle)
	job = cron.NewChain(cron.SkipIfStillRunning(cronLog)).Then(job)
	_, err := cronJob.AddJob(config.ProcessingInterval(), job)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	// Cleanup job
	job = cron.FuncJob(cleanupHandler.Handle)
	job = cron.NewChain(cron.SkipIfStillRunning(cronLog)).Then(job)
	_, err = cronJob.AddJob(config.CleanupInterval(), job)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	cronJob.Start()

	//Run forever
	select {}
}