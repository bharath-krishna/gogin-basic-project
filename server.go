package main

import (
	"os"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Server struct {
	authConfig oauth2.Config
	provider   *oidc.Provider
	gclient    *Client
	config     *Config
	app        *gin.Engine
	logger     *zap.Logger
}

func createNewServer(config *Config) (*Server, error) {
	logger, err := createLogger(config)
	if err != nil {
		return nil, err
	}

	logger.Info("Starting the service", zap.String("prog", prog), zap.String("version", version))

	provider, err := oidc.NewProvider(oauth2.NoContext, config.AuthHost+"/auth/realms/demo")
	if err != nil {
		panic(err)
	}
	config.provider = provider

	// For three legged authentication flow
	authConfig := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  os.Getenv("CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{"profile", "email"},
	}

	gclient, err := getNewGraphClient(logger, config)
	if err != nil {
		return nil, err
	}
	s := &Server{
		authConfig: authConfig,
		gclient:    gclient,
		logger:     logger,
		provider:   config.provider,
		config:     config,
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
