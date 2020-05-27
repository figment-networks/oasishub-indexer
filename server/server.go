package server

import (
	"github.com/figment-networks/oasishub-indexer/usecase"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/gin-gonic/gin"
)

// Server handles HTTP requests
type Server struct {
	engine *gin.Engine

	handlers *usecase.HttpHandlers
}

// New returns a new server instance
func New(handlers *usecase.HttpHandlers) *Server {
	app := &Server{
		engine: gin.Default(),
		handlers: handlers,
	}
	return app.init()
}

func (a *Server) init() *Server {
	logger.Info("initializing server...")

	a.engine.GET("/health", a.handlers.Health.Handle)
	a.engine.GET("/blocks", a.handlers.GetBlockByHeight.Handle)
	a.engine.GET("/block_times/:limit", a.handlers.GetBlockTimes.Handle)
	a.engine.GET("/block_times_interval", a.handlers.GetBlockTimesForInterval.Handle)
	a.engine.GET("/transactions", a.handlers.GetTransactionsByHeight.Handle)
	a.engine.GET("/validators/by_entity_uid", a.handlers.GetValidatorByEntityUid.Handle)
	a.engine.GET("/validators", a.handlers.GetValidatorsByHeight.Handle)
	a.engine.GET("/validators/for_min_height/:height", a.handlers.GetValidatorsForMinHeight.Handle)
	a.engine.GET("/validators/shares_interval", a.handlers.GetValidatorShares.Handle)
	a.engine.GET("/validators/voting_power_interval", a.handlers.GetValidatorVotingPower.Handle)
	a.engine.GET("/validators/uptime_interval", a.handlers.GetValidatorUptime.Handle)
	a.engine.GET("/validators/total_shares_interval", a.handlers.GetSharesForAllValidators.Handle)
	a.engine.GET("/validators/total_voting_power_interval", a.handlers.GetVotingPowerForAllValidators.Handle)
	a.engine.GET("/staking", a.handlers.GetStakingDetailsByHeight.Handle)
	a.engine.GET("/delegations", a.handlers.GetDebondingDelegationsByHeight.Handle)
	a.engine.GET("/debonding_delegations", a.handlers.GetDebondingDelegationsByHeight.Handle)
	a.engine.GET("/accounts", a.handlers.GetAccountByPublicKey.Handle)
	a.engine.GET("/current_height", a.handlers.GetMostRecentHeight.Handle)
	return a
}

func (a *Server) Start(listenAdd string) error {
	logger.Info("starting server...", logger.Field("app", "server"))
	return a.engine.Run(listenAdd)
}