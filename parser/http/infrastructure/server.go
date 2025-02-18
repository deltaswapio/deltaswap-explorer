package infrastructure

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/deltaswapio/deltaswap-explorer/parser/http/vaa"
	"github.com/deltaswapio/deltaswap-explorer/parser/internal/sqs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Server struct {
	app    *fiber.App
	port   string
	logger *zap.Logger
}

func NewServer(logger *zap.Logger, port string, pprofEnabled bool, isQueueConsumer bool, consumer *sqs.Consumer,
	db *mongo.Database, vaaController *vaa.Controller) *Server {
	repository := NewRepository(db, logger)
	service := NewService(repository, consumer, isQueueConsumer, logger)
	ctrl := NewController(service, logger)
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// config use of middlware.
	prometheus := fiberprometheus.New("deltaswapscan-pipeline")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	if pprofEnabled {
		app.Use(pprof.New())
	}

	api := app.Group("/api")
	api.Get("/health", ctrl.HealthCheck)
	api.Get("/ready", ctrl.ReadyCheck)

	api.Post("/vaa/parse", vaaController.Parse)

	return &Server{
		app:    app,
		port:   port,
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
