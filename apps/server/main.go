package main

import (
	"fmt"
	"github.com/figment-networks/oasishub/apps/shared"
	"github.com/figment-networks/oasishub/config"
	"github.com/figment-networks/oasishub/repos/blockseqrepo"
	"github.com/figment-networks/oasishub/repos/syncablerepo"
	"github.com/figment-networks/oasishub/usecases/block/getblockbyheight"
	"github.com/figment-networks/oasishub/usecases/block/getblocktimes"
	"github.com/figment-networks/oasishub/usecases/block/getblocktimesforinterval"
	"github.com/figment-networks/oasishub/usecases/ping"
	"github.com/figment-networks/oasishub/utils/log"
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

	blockDbRepo := blockseqrepo.NewDbRepo(db.Client())

	//USE CASES
	getBlockByHeight := getblockbyheight.NewUseCase(syncableDbRepo, syncableProxyRepo, blockDbRepo)
	uc := getblocktimes.NewUseCase(blockDbRepo)
	uc2 := getblocktimesforinterval.NewUseCase(blockDbRepo)

	// HANDLERS
	pingHandler := ping.NewHttpHandler()
	getBlockByHeightHandler := getblockbyheight.NewHttpHandler(getBlockByHeight)
	getAvgBlockTimesForRecent := getblocktimes.NewHttpHandler(uc)
	getAvgBlockTimesForInterval := getblocktimesforinterval.NewHttpHandler(uc2)

	// ADD ROUTES
	router = gin.Default()
	router.GET("/ping", pingHandler.Handle)
	router.GET("/blocks/:height", getBlockByHeightHandler.Handle)
	router.GET("/block-times/:limit", getAvgBlockTimesForRecent.Handle)
	router.GET("/block-times-interval/:interval", getAvgBlockTimesForInterval.Handle)

	log.Info(fmt.Sprintf("Starting application on port %s", config.AppPort()))

	// START SERVER
	if err := router.Run(fmt.Sprintf(":%s", config.AppPort())); err != nil {
		panic(err)
	}
}
