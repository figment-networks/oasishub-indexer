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
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecases/syncable/startpipeline"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	BatchSize = "batchSize"
)

var (
	rootCmd   *cobra.Command
	batchSize int64
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

	// HANDLERS
	startPipelineCliHandler := startpipeline.NewCliHandler(startPipelineUseCase)

	// CLI COMMANDS
	rootCmd = setupRootCmd()
	pipelineCmd := setupPipelineCmd(startPipelineCliHandler)

	rootCmd.AddCommand(pipelineCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		panic(err)
	}
}

/*************** Private ***************/

func setupRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cli",
		Short: "Short description",
		Long: `Longer description.. 
            feel free to use a few lines here.
            `,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); err != nil {
				log.Error(err, log.Field("type", "cli"))
			}
		},
	}
}

func setupPipelineCmd(handler types.CliHandler) *cobra.Command {
	pipelineCmd := &cobra.Command{
		Use:   "pipeline [command]",
		Short: "Run pipeline commands",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); err != nil {
				log.Error(err, log.Field("type", "cli"))
			}
		},
	}

	startPipelineCmd := &cobra.Command{
		Use:   "start",
		Short: "Start one off processing pipeline",
		Args:  cobra.MaximumNArgs(1),
		Run:   handler.Handle,
	}
	rootCmd.PersistentFlags().Int64Var(&batchSize, BatchSize, config.PipelineBatchSize(), "batch size")
	viper.BindPFlag(BatchSize, rootCmd.PersistentFlags().Lookup(BatchSize))
	pipelineCmd.AddCommand(startPipelineCmd)
	return pipelineCmd
}
