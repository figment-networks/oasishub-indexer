package server

import (
	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/metrics/prometheusmetrics"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/usecase"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/gin-gonic/gin"
)

// Server handles HTTP requests
type Server struct {
	cfg      *config.Config
	handlers *usecase.HttpHandlers

	engine *gin.Engine
}

// New returns a new server instance
func New(cfg *config.Config, handlers *usecase.HttpHandlers) *Server {
	app := &Server{
		cfg:      cfg,
		engine:   gin.Default(),
		handlers: handlers,
	}
	return app.init()
}

// Index starts the server
func (s *Server) Start(listenAdd string) error {
	logger.Info("starting server...", logger.Field("app", "server"))

	prom := prometheusmetrics.New()
	err := metrics.AddEngine(prom)
	if err != nil {
		logger.Error(err)
	}
	err = metrics.Hotload(prom.Name())
	if err != nil {
		logger.Error(err)
	}
	s.engine.GET(s.cfg.MetricServerUrl, gin.WrapH(metrics.Handler()))

	return s.engine.Run(listenAdd)
}

// init initializes the server
func (s *Server) init() *Server {
	logger.Info("initializing server...", logger.Field("app", "server"))

	s.setupMiddleware()
	s.setupRoutes()

	return s
}
