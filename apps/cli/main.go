package main

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasishub-indexer/apps/shared"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/chainrepo"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/entityaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/reportrepo"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecases/pipeline/startpipeline"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
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
	defer errors.RecoverError()

	// CLIENTS
	proxy := shared.NewProxyClient()
	defer proxy.Client().Close()

	db := shared.NewDbClient()
	defer db.Client().Close()

	// REPOSITORIES
	chainProxyRepo := chainrepo.NewProxyRepo(chainpb.NewChainServiceClient(proxy.Client()))
	syncableDbRepo := syncablerepo.NewDbRepo(db.Client())
	syncableProxyRepo := syncablerepo.NewProxyRepo(
		blockpb.NewBlockServiceClient(proxy.Client()),
		statepb.NewStateServiceClient(proxy.Client()),
		transactionpb.NewTransactionServiceClient(proxy.Client()),
		validatorpb.NewValidatorServiceClient(proxy.Client()),
	)
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
		chainProxyRepo,
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
