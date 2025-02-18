package infrastructure

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/deltaswapio/deltaswap-explorer/pipeline/healthcheck"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"go.uber.org/zap"
)

type Server struct {
	app    *fiber.App
	port   string
	logger *zap.Logger
}

func NewServer(logger *zap.Logger, port string, pprofEnabled bool, checks ...healthcheck.Check) *Server {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// config use of middlware.
	prometheus := fiberprometheus.New("deltaswapscan-pipeline")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	if pprofEnabled {
		app.Use(pprof.New())
	}

	ctrl := NewController(checks, logger)
	api := app.Group("/api")
	api.Get("/health", ctrl.HealthCheck)
	api.Get("/ready", ctrl.ReadyCheck)

	return &Server{
		app:    app,
		port:   port,
		logger: logger,
	}
}

// Start listen serves HTTP requests from addr.
func (s *Server) Start() {
	addr := ":" + s.port
	s.logger.Info("Listening on " + addr)
	go func() {
		s.app.Listen(addr)
	}()
}

// Stop gracefull server.
func (s *Server) Stop() {
	_ = s.app.Shutdown()
}
