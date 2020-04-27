package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	gclient *Client
	config  *Config
	app     *gin.Engine
	logger  *zap.Logger
}

func createNewServer(config *Config) (*Server, error) {
	logger, err := createLogger(config)
	if err != nil {
		return nil, err
	}

	logger.Info("Starting the service", zap.String("prog", prog), zap.String("version", version))

	gclient, err := getNewGraphClient(logger, config)
	if err != nil {
		return nil, err
	}
	s := &Server{
		gclient: gclient,
		logger:  logger,
		config:  config,
	}

	s.app = s.setupRoutes()
	return s, nil
}

func (s *Server) setupRoutes() *gin.Engine {
	r := gin.New()
	r.Use(s.LoggingMiddleware())
	r.Use(gin.Recovery())
	r.Use(s.AllowOriginRequests())
	s.Routes(r)
	return r
}

// Run runs the server
func (s *Server) Run() {
	s.app.Run(s.config.ServerAddr)
}

func createLogger(config *Config) (*zap.Logger, error) {
	c := zap.NewProductionConfig()
	c.DisableCaller = true
	// c.Encoding = "console"

	if config.Verbose {
		c.DisableCaller = false
		c.Development = true
		c.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	return c.Build()
}
