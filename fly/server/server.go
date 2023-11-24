package server

import (
	"fmt"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/deltaswapio/deltaswap-explorer/common/client/alert"
	"github.com/deltaswapio/deltaswap-explorer/fly/internal/health"
	"github.com/deltaswapio/deltaswap-explorer/fly/internal/sqs"
	"github.com/deltaswapio/deltaswap-explorer/fly/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"go.uber.org/zap"
)

type Server struct {
	app    *fiber.App
	port   string
	logger *zap.Logger
}

func NewServer(port uint, phylaxCheck *health.PhylaxCheck, logger *zap.Logger, repository *storage.Repository, consumer *sqs.Consumer, isLocal, pprofEnabled bool, alertClient alert.AlertClient) *Server {
	ctrl := NewController(phylaxCheck, repository, consumer, isLocal, alertClient, logger)
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// Configure middleware
	prometheus := fiberprometheus.New("wormscan-fly")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// config use of middlware.
	if pprofEnabled {
		app.Use(pprof.New())
	}
	api := app.Group("/api")
	api.Get("/health", ctrl.HealthCheck)
	api.Get("/ready", ctrl.ReadyCheck)
	return &Server{
		app:    app,
		port:   fmt.Sprintf("%d", port),
		logger: logger,
	}
}

// Start listen serves HTTP requests from addr.
func (s *Server) Start() {
	go func() {
		s.app.Listen(":" + s.port)
	}()
}

// Stop gracefull server.
func (s *Server) Stop() {
	_ = s.app.Shutdown()
}
