package main

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/apps/shared"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/entityaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/usecases/account/getaccountbypublickey"
	"github.com/figment-networks/oasishub-indexer/usecases/block/getblockbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/block/getblocktimes"
	"github.com/figment-networks/oasishub-indexer/usecases/block/getblocktimesforinterval"
	"github.com/figment-networks/oasishub-indexer/usecases/delegation/getdebondingdelegationsbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/delegation/getdelegationsbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/ping"
	"github.com/figment-networks/oasishub-indexer/usecases/staking/getstakingbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/transaction/gettransactionsbyheight"
	"github.com/figment-networks/oasishub-indexer/usecases/validator/getentitybyentityuid"
	"github.com/figment-networks/oasishub-indexer/usecases/validator/getvalidatorsbyheight"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func main() {
	// CLIENTS
	node := shared.NewNodeClient()
	db := shared.NewDbClient()

	// REPOSITORIES
	syncableDbRepo := syncablerepo.NewDbRepo(db.Client())
	syncableProxyRepo := syncablerepo.NewProxyRepo(node)

	blockSeqDbRepo := blockseqrepo.NewDbRepo(db.Client())
	transactionSeqDbRepo := transactionseqrepo.NewDbRepo(db.Client())
	validatorSeqDbRepo := validatorseqrepo.NewDbRepo(db.Client())
	stakingSeqDbRepo := stakingseqrepo.NewDbRepo(db.Client())
	delegationSeqDbRepo := delegationseqrepo.NewDbRepo(db.Client())
	debondingDelegationSeqDbRepo := debondingdelegationseqrepo.NewDbRepo(db.Client())
	accountAggDbRepo := accountaggrepo.NewDbRepo(db.Client())
	entityAggDbRepo := entityaggrepo.NewDbRepo(db.Client())

	//USE CASES
	getBlockByHeight := getblockbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, blockSeqDbRepo, validatorSeqDbRepo, transactionSeqDbRepo)
	getBlockTimes := getblocktimes.NewUseCase(blockSeqDbRepo)
	getBlockTimesForInterval := getblocktimesforinterval.NewUseCase(blockSeqDbRepo)
	getTransactionsByHeight := gettransactionsbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, transactionSeqDbRepo)
	getValidatorsByHeight := getvalidatorsbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, validatorSeqDbRepo, delegationSeqDbRepo)
	getStakingByHeight := getstakingbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, stakingSeqDbRepo)
	getDelegationsByHeight := getdelegationsbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, delegationSeqDbRepo)
	getDebondingDelegationsByHeight := getdebondingdelegationsbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, debondingDelegationSeqDbRepo)
	getAccountByPublicKey := getaccountbypublickey.NewUseCase(syncableDbRepo, syncableProxyRepo, accountAggDbRepo, delegationSeqDbRepo, debondingDelegationSeqDbRepo)
	getEntityByEntityUID := getentitybyentityuid.NewUseCase(syncableDbRepo, syncableProxyRepo, entityAggDbRepo, validatorSeqDbRepo, delegationSeqDbRepo, debondingDelegationSeqDbRepo)

	// HANDLERS
	pingHandler := ping.NewHttpHandler()
	getBlockByHeightHandler := getblockbyheight.NewHttpHandler(getBlockByHeight)
	getAvgBlockTimesForRecentHandler := getblocktimes.NewHttpHandler(getBlockTimes)
	getAvgBlockTimesForIntervalHandler := getblocktimesforinterval.NewHttpHandler(getBlockTimesForInterval)
	getTransactionsByHeightHandler := gettransactionsbyheight.NewHttpHandler(getTransactionsByHeight)
	getValidatorsByHeightHandler := getvalidatorsbyheight.NewHttpHandler(getValidatorsByHeight)
	getStakingByHeightHandler := getstakingbyheight.NewHttpHandler(getStakingByHeight)
	getDelegationsByHeightHandler := getdelegationsbyheight.NewHttpHandler(getDelegationsByHeight)
	getDebondingDelegationsByHeightHandler := getdebondingdelegationsbyheight.NewHttpHandler(getDebondingDelegationsByHeight)
	getAccountByPublicKeyHandler := getaccountbypublickey.NewHttpHandler(getAccountByPublicKey)
	getEntityByEntityUIDHandler := getentitybyentityuid.NewHttpHandler(getEntityByEntityUID)

	// ADD ROUTES
	router = gin.Default()
	router.GET("/ping", pingHandler.Handle)
	router.GET("/blocks", getBlockByHeightHandler.Handle)
	router.GET("/block_times/:limit", getAvgBlockTimesForRecentHandler.Handle)
	router.GET("/block_times_interval/:interval", getAvgBlockTimesForIntervalHandler.Handle)
	router.GET("/transactions", getTransactionsByHeightHandler.Handle)
	router.GET("/entities", getEntityByEntityUIDHandler.Handle)
	router.GET("/validators", getValidatorsByHeightHandler.Handle)
	router.GET("/staking", getStakingByHeightHandler.Handle)
	router.GET("/delegations", getDelegationsByHeightHandler.Handle)
	router.GET("/debonding_delegations", getDebondingDelegationsByHeightHandler.Handle)
	router.GET("/accounts", getAccountByPublicKeyHandler.Handle)

	log.Info(fmt.Sprintf("Starting application on port %s", config.AppPort()))

	// START SERVER
	if err := router.Run(fmt.Sprintf(":%s", config.AppPort())); err != nil {
		panic(err)
	}
}
