package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/deltaswapio/deltaswap-explorer/common/client/sqs"
	"github.com/deltaswapio/deltaswap-explorer/common/dbutil"
	"github.com/deltaswapio/deltaswap-explorer/common/health"
	"github.com/deltaswapio/deltaswap-explorer/common/logger"
	"github.com/deltaswapio/deltaswap-explorer/txtracker/chains"
	"github.com/deltaswapio/deltaswap-explorer/txtracker/config"
	"github.com/deltaswapio/deltaswap-explorer/txtracker/consumer"
	"github.com/deltaswapio/deltaswap-explorer/txtracker/http/infrastructure"
	"github.com/deltaswapio/deltaswap-explorer/txtracker/http/vaa"
	"github.com/deltaswapio/deltaswap-explorer/txtracker/internal/metrics"
	"github.com/deltaswapio/deltaswap-explorer/txtracker/queue"
	"go.uber.org/zap"
)

func main() {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())

	// load config
	cfg, err := config.LoadFromEnv[config.ServiceSettings]()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	// initialize metrics
	metrics := newMetrics(cfg)

	// build logger
	logger := logger.New("deltaswap-explorer-tx-tracker", logger.WithLevel(cfg.LogLevel))

	logger.Info("Starting deltaswap-explorer-tx-tracker ...")

	// initialize rate limiters
	chains.Initialize(&cfg.RpcProviderSettings)

	// initialize the database client
	db, err := dbutil.Connect(rootCtx, logger, cfg.MongodbUri, cfg.MongodbDatabase, false)
	if err != nil {
		log.Fatal("Failed to initialize MongoDB client: ", err)
	}

	// create repositories
	repository := consumer.NewRepository(logger, db.Database)
	vaaRepository := vaa.NewRepository(db.Database, logger)

	// create controller
	vaaController := vaa.NewController(vaaRepository, repository, &cfg.RpcProviderSettings, cfg.P2pNetwork, logger)

	// start serving /health and /ready endpoints
	healthChecks, err := makeHealthChecks(rootCtx, cfg, db.Database)
	if err != nil {
		logger.Fatal("Failed to create health checks", zap.Error(err))
	}
	server := infrastructure.NewServer(logger, cfg.MonitoringPort, cfg.PprofEnabled, vaaController, healthChecks...)
	server.Start()

	// create and start a consumer.
	vaaConsumeFunc := newVAAConsumeFunc(rootCtx, cfg, metrics, logger)
	consumer := consumer.New(vaaConsumeFunc, &cfg.RpcProviderSettings, rootCtx, logger, repository, metrics, cfg.P2pNetwork)
	consumer.Start(rootCtx)

	logger.Info("Started deltaswap-explorer-tx-tracker")

	// Waiting for signal
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-rootCtx.Done():
		logger.Warn("Terminating with root context cancelled.")
	case signal := <-sigterm:
		logger.Info("Terminating with signal.", zap.String("signal", signal.String()))
	}

	// graceful shutdown
	logger.Info("Cancelling root context...")
	rootCtxCancel()

	logger.Info("Closing Http server...")
	server.Stop()

	logger.Info("Closing MongoDB connection...")
	db.DisconnectWithTimeout(10 * time.Second)

	logger.Info("Terminated deltaswap-explorer-tx-tracker")
}

func newVAAConsumeFunc(
	ctx context.Context,
	cfg *config.ServiceSettings,
	metrics metrics.Metrics,
	logger *zap.Logger,
) queue.VAAConsumeFunc {

	sqsConsumer, err := newSqsConsumer(ctx, cfg)
	if err != nil {
		logger.Fatal("failed to create sqs consumer", zap.Error(err))
	}

	vaaQueue := queue.NewVaaSqs(sqsConsumer, metrics, logger)
	return vaaQueue.Consume
}

func newSqsConsumer(ctx context.Context, cfg *config.ServiceSettings) (*sqs.Consumer, error) {

	awsconfig, err := newAwsConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	consumer, err := sqs.NewConsumer(
		awsconfig,
		cfg.PipelineSqsUrl,
		sqs.WithMaxMessages(10),
		sqs.WithVisibilityTimeout(4*60),
	)
	return consumer, err
}

func newAwsConfig(ctx context.Context, cfg *config.ServiceSettings) (aws.Config, error) {

	region := cfg.AwsRegion

	if cfg.AwsAccessKeyID != "" && cfg.AwsSecretAccessKey != "" {

		credentials := credentials.NewStaticCredentialsProvider(cfg.AwsAccessKeyID, cfg.AwsSecretAccessKey, "")

		customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if cfg.AwsEndpoint != "" {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           cfg.AwsEndpoint,
					SigningRegion: region,
				}, nil
			}

			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		awsCfg, err := awsconfig.LoadDefaultConfig(
			ctx,
			awsconfig.WithRegion(region),
			awsconfig.WithEndpointResolver(customResolver),
			awsconfig.WithCredentialsProvider(credentials),
		)
		return awsCfg, err
	}
	return awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
}

func makeHealthChecks(
	ctx context.Context,
	config *config.ServiceSettings,
	db *mongo.Database,
) ([]health.Check, error) {

	awsConfig, err := newAwsConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	plugins := []health.Check{
		health.SQS(awsConfig, config.PipelineSqsUrl),
		health.Mongo(db),
	}

	return plugins, nil
}

func newMetrics(cfg *config.ServiceSettings) metrics.Metrics {
	if !cfg.MetricsEnabled {
		return metrics.NewDummyMetrics()
	}
	return metrics.NewPrometheusMetrics(cfg.Environment)
}
