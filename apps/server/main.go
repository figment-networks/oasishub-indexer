package main

import (
	"fmt"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasishub-indexer/apps/shared"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatoraggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/usecases/account/getaccountbypublickey"
	"github.com/figment-networks/oasishub-indexer/usecases/block/getblockbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/block/getblocktimes"
	"github.com/figment-networks/oasishub-indexer/usecases/block/getblocktimesforinterval"
	"github.com/figment-networks/oasishub-indexer/usecases/debondingdelegation/getdebondingdelegationsbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/delegation/getdelegationsbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/ping"
	"github.com/figment-networks/oasishub-indexer/usecases/staking/getstakingbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/syncable/getmostrecentheight"
	"github.com/figment-networks/oasishub-indexer/usecases/transaction/gettransactionsbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/validator/gettotalsharesforinterval"
	"github.com/figment-networks/oasishub-indexer/usecases/validator/gettotalvotingpowerforinterval"
	"github.com/figment-networks/oasishub-indexer/usecases/validator/getvalidatorbyentityuid"
	"github.com/figment-networks/oasishub-indexer/usecases/validator/getvalidatorsbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/validator/getvalidatorsforminheight"
	"github.com/figment-networks/oasishub-indexer/usecases/validator/getvalidatorsharesforinterval"
	"github.com/figment-networks/oasishub-indexer/usecases/validator/getvalidatoruptimeforinterval"
	"github.com/figment-networks/oasishub-indexer/usecases/validator/getvalidatorvotingpowerforinterval"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func main() {
	defer errors.RecoverError()

	// CLIENTS
	proxy := shared.NewProxyClient()
	defer proxy.Client().Close()

	db := shared.NewDbClient()
	defer db.Client().Close()

	// REPOSITORIES
	syncableDbRepo := syncablerepo.NewDbRepo(db.Client())
	syncableProxyRepo := syncablerepo.NewProxyRepo(
		blockpb.NewBlockServiceClient(proxy.Client()),
		statepb.NewStateServiceClient(proxy.Client()),
		transactionpb.NewTransactionServiceClient(proxy.Client()),
		validatorpb.NewValidatorServiceClient(proxy.Client()),
	)

	blockSeqDbRepo := blockseqrepo.NewDbRepo(db.Client())
	transactionSeqDbRepo := transactionseqrepo.NewDbRepo(db.Client())
	validatorSeqDbRepo := validatorseqrepo.NewDbRepo(db.Client())
	stakingSeqDbRepo := stakingseqrepo.NewDbRepo(db.Client())
	delegationSeqDbRepo := delegationseqrepo.NewDbRepo(db.Client())
	debondingDelegationSeqDbRepo := debondingdelegationseqrepo.NewDbRepo(db.Client())
	accountAggDbRepo := accountaggrepo.NewDbRepo(db.Client())
	validatorAggDbRepo := validatoraggrepo.NewDbRepo(db.Client())

	//USE CASES
	getBlockByHeight := getblockbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, blockSeqDbRepo, validatorSeqDbRepo, transactionSeqDbRepo)
	getBlockTimes := getblocktimes.NewUseCase(blockSeqDbRepo)
	getBlockTimesForInterval := getblocktimesforinterval.NewUseCase(blockSeqDbRepo)
	getTransactionsByHeight := gettransactionsbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, transactionSeqDbRepo)
	getValidatorsByHeight := getvalidatorsbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, validatorSeqDbRepo, delegationSeqDbRepo)
	getValidatorSharesByInterval := getvalidatorsharesforinterval.NewUseCase(syncableDbRepo, syncableProxyRepo, validatorSeqDbRepo)
	getValidatorVotingPowerByInterval := getvalidatorvotingpowerforinterval.NewUseCase(syncableDbRepo, syncableProxyRepo, validatorSeqDbRepo)
	getValidatorUptimeByInterval := getvalidatoruptimeforinterval.NewUseCase(syncableDbRepo, syncableProxyRepo, validatorSeqDbRepo)
	getTotalSharesByInterval := gettotalsharesforinterval.NewUseCase(syncableDbRepo, syncableProxyRepo, validatorSeqDbRepo)
	getTotalVotingPowerByInterval := gettotalvotingpowerforinterval.NewUseCase(validatorSeqDbRepo)
	getStakingByHeight := getstakingbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, stakingSeqDbRepo)
	getDelegationsByHeight := getdelegationsbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, delegationSeqDbRepo)
	getDebondingDelegationsByHeight := getdebondingdelegationsbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, debondingDelegationSeqDbRepo)
	getAccountByPublicKey := getaccountbypublickey.NewUseCase(syncableDbRepo, syncableProxyRepo, accountAggDbRepo, delegationSeqDbRepo, debondingDelegationSeqDbRepo)
	getValidatorByEntityUID := getvalidatorbyentityuid.NewUseCase(syncableDbRepo, syncableProxyRepo, validatorAggDbRepo, validatorSeqDbRepo, delegationSeqDbRepo, debondingDelegationSeqDbRepo)
	getValidatorsForMinHeight := getvalidatorsforminheight.NewUseCase(syncableDbRepo, syncableProxyRepo, validatorAggDbRepo)
	getMostRecentHeight := getmostrecentheight.NewUseCase(syncableDbRepo)

	// HANDLERS
	pingHandler := ping.NewHttpHandler()
	getBlockByHeightHandler := getblockbyheight.NewHttpHandler(getBlockByHeight)
	getAvgBlockTimesForRecentHandler := getblocktimes.NewHttpHandler(getBlockTimes)
	getAvgBlockTimesForIntervalHandler := getblocktimesforinterval.NewHttpHandler(getBlockTimesForInterval)
	getTransactionsByHeightHandler := gettransactionsbyheight.NewHttpHandler(getTransactionsByHeight)
	getValidatorsByHeightHandler := getvalidatorsbyheight.NewHttpHandler(getValidatorsByHeight)
	getValidatorSharesByIntervalHandler := getvalidatorsharesforinterval.NewHttpHandler(getValidatorSharesByInterval)
	getValidatorVotingPowerByIntervalHandler := getvalidatorvotingpowerforinterval.NewHttpHandler(getValidatorVotingPowerByInterval)
	getValidatorUptimeByIntervalHandler := getvalidatoruptimeforinterval.NewHttpHandler(getValidatorUptimeByInterval)
	getTotalSharesByIntervalHandler := gettotalsharesforinterval.NewHttpHandler(getTotalSharesByInterval)
	getTotalVotingPowerByIntervalHandler := gettotalvotingpowerforinterval.NewHttpHandler(getTotalVotingPowerByInterval)
	getStakingByHeightHandler := getstakingbyheight.NewHttpHandler(getStakingByHeight)
	getDelegationsByHeightHandler := getdelegationsbyheight.NewHttpHandler(getDelegationsByHeight)
	getDebondingDelegationsByHeightHandler := getdebondingdelegationsbyheight.NewHttpHandler(getDebondingDelegationsByHeight)
	getAccountByPublicKeyHandler := getaccountbypublickey.NewHttpHandler(getAccountByPublicKey)
	getValidatorByEntityUIDHandler := getvalidatorbyentityuid.NewHttpHandler(getValidatorByEntityUID)
	getValidatorsForMinHeightHandler := getvalidatorsforminheight.NewHttpHandler(getValidatorsForMinHeight)
	getMostRecentHeightHandler := getmostrecentheight.NewHttpHandler(getMostRecentHeight)

	// ADD ROUTES
	router = gin.Default()
	router.GET("/ping", pingHandler.Handle)
	router.GET("/blocks", getBlockByHeightHandler.Handle)
	router.GET("/block_times/:limit", getAvgBlockTimesForRecentHandler.Handle)
	router.GET("/block_times_interval", getAvgBlockTimesForIntervalHandler.Handle)
	router.GET("/transactions", getTransactionsByHeightHandler.Handle)
	router.GET("/validators/by_entity_uid", getValidatorByEntityUIDHandler.Handle)
	router.GET("/validators", getValidatorsByHeightHandler.Handle)
	router.GET("/validators/for_min_height/:height", getValidatorsForMinHeightHandler.Handle)
	router.GET("/validators/shares_interval", getValidatorSharesByIntervalHandler.Handle)
	router.GET("/validators/voting_power_interval", getValidatorVotingPowerByIntervalHandler.Handle)
	router.GET("/validators/uptime_interval", getValidatorUptimeByIntervalHandler.Handle)
	router.GET("/validators/total_shares_interval", getTotalSharesByIntervalHandler.Handle)
	router.GET("/validators/total_voting_power_interval", getTotalVotingPowerByIntervalHandler.Handle)
	router.GET("/staking", getStakingByHeightHandler.Handle)
	router.GET("/delegations", getDelegationsByHeightHandler.Handle)
	router.GET("/debonding_delegations", getDebondingDelegationsByHeightHandler.Handle)
	router.GET("/accounts", getAccountByPublicKeyHandler.Handle)
	router.GET("/current_height", getMostRecentHeightHandler.Handle)

	log.Info(fmt.Sprintf("Starting server on port %s", config.AppPort()))

	// START SERVER
	if err := router.Run(fmt.Sprintf(":%s", config.AppPort())); err != nil {
		panic(err)
	}
}
