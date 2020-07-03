package server

import (
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

	//s.cfg.ServerMetricAddr, s.cfg.MetricServerUrl
	return s.engine.Run(listenAdd)
}

// init initializes the server
func (s *Server) init() *Server {
	logger.Info("initializing server...", logger.Field("app", "server"))

	s.setupMiddleware()
	s.setupRoutes()

	return s
}
