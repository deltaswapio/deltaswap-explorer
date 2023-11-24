package service

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
	"github.com/deltaswapio/deltaswap-explorer/common/client/alert"
	vaaPayloadParser "github.com/deltaswapio/deltaswap-explorer/common/client/parser"
	"github.com/deltaswapio/deltaswap-explorer/common/dbutil"
	"github.com/deltaswapio/deltaswap-explorer/common/logger"
	"github.com/deltaswapio/deltaswap-explorer/parser/config"
	"github.com/deltaswapio/deltaswap-explorer/parser/consumer"
	"github.com/deltaswapio/deltaswap-explorer/parser/http/infrastructure"
	"github.com/deltaswapio/deltaswap-explorer/parser/http/vaa"
	parserAlert "github.com/deltaswapio/deltaswap-explorer/parser/internal/alert"
	"github.com/deltaswapio/deltaswap-explorer/parser/internal/metrics"
	"github.com/deltaswapio/deltaswap-explorer/parser/internal/sqs"
	"github.com/deltaswapio/deltaswap-explorer/parser/migration"
	"github.com/deltaswapio/deltaswap-explorer/parser/parser"
	"github.com/deltaswapio/deltaswap-explorer/parser/processor"
	"github.com/deltaswapio/deltaswap-explorer/parser/queue"
	"go.uber.org/zap"
)

type exitCode int

func handleExit() {
	if r := recover(); r != nil {
		if e, ok := r.(exitCode); ok {
			os.Exit(int(e))
		}
		panic(r) // not an Exit, bubble up
	}
}

func Run() {

	defer handleExit()
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())

	config, err := config.New(rootCtx)
	if err != nil {
		log.Fatal("Error creating config", err)
	}

	logger := logger.New("deltaswap-explorer-parser", logger.WithLevel(config.LogLevel))

	logger.Info("Starting deltaswap-explorer-parser ...")

	// setup DB connection
	db, err := dbutil.Connect(rootCtx, logger, config.MongoURI, config.MongoDatabase, false)
	if err != nil {
		logger.Fatal("failed to connect MongoDB", zap.Error(err))
	}

	// run the database migration.
	err = migration.Run(db.Database)
	if err != nil {
		logger.Fatal("error running migration", zap.Error(err))
	}

	// get alert client.
	alertClient, err := newAlertClient(config)
	if err != nil {
		logger.Fatal("failed to create alert client", zap.Error(err))
	}

	// create a metrics
	metrics := newMetrics(config)

	// create a parserVAAAPIClient
	parserVAAAPIClient, err := vaaPayloadParser.NewParserVAAAPIClient(config.VaaPayloadParserTimeout,
		config.VaaPayloadParserURL, logger)
	if err != nil {
		logger.Fatal("failed to create parse vaa api client")
	}

	// get consumer function.
	sqsConsumer, vaaConsumeFunc := newVAAConsume(rootCtx, config, metrics, logger)
	repository := parser.NewRepository(db.Database, logger)

	//create a processor
	processor := processor.New(parserVAAAPIClient, repository, alertClient, metrics, logger)

	// create and start a consumer
	consumer := consumer.New(vaaConsumeFunc, processor.Process, metrics, logger)
	consumer.Start(rootCtx)

	vaaRepository := vaa.NewRepository(db.Database, logger)
	vaaController := vaa.NewController(vaaRepository, processor.Process, logger)
	server := infrastructure.NewServer(logger, config.Port, config.PprofEnabled, config.IsQueueConsumer(), sqsConsumer, db.Database, vaaController)
	server.Start()

	logger.Info("Started deltaswap-explorer-parser")

	// Waiting for signal
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-rootCtx.Done():
		logger.Warn("Terminating with root context cancelled.")
	case signal := <-sigterm:
		logger.Info("Terminating with signal.", zap.String("signal", signal.String()))
	}

	logger.Info("root context cancelled, exiting...")
	rootCtxCancel()

	logger.Info("closing MongoDB connection...")
	db.DisconnectWithTimeout(10 * time.Second)

	logger.Info("Closing Http server ...")
	server.Stop()

	logger.Info("Finished deltaswap-explorer-parser")
}

// Creates a new AWS config depending on whether the execution is local (localstack) or not (AWS)
func newAwsConfig(appCtx context.Context, cfg *config.ServiceConfiguration) (aws.Config, error) {
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

		awsCfg, err := awsconfig.LoadDefaultConfig(appCtx,
			awsconfig.WithRegion(region),
			awsconfig.WithEndpointResolver(customResolver),
			awsconfig.WithCredentialsProvider(credentials),
		)
		return awsCfg, err
	}
	return awsconfig.LoadDefaultConfig(appCtx, awsconfig.WithRegion(region))
}

func newVAAConsume(appCtx context.Context, config *config.ServiceConfiguration, metrics metrics.Metrics, logger *zap.Logger) (*sqs.Consumer, queue.VAAConsumeFunc) {
	sqsConsumer, err := newSQSConsumer(appCtx, config)
	if err != nil {
		logger.Fatal("failed to create sqs consumer", zap.Error(err))
	}

	filterConsumeFunc := newFilterFunc(config)
	vaaQueue := queue.NewVAASQS(sqsConsumer, filterConsumeFunc, metrics, logger)
	return sqsConsumer, vaaQueue.Consume
}

// Create a new SQS consumer.
func newSQSConsumer(appCtx context.Context, config *config.ServiceConfiguration) (*sqs.Consumer, error) {
	awsconfig, err := newAwsConfig(appCtx, config)
	if err != nil {
		return nil, err
	}

	return sqs.NewConsumer(awsconfig, config.SQSUrl,
		sqs.WithMaxMessages(10),
		sqs.WithVisibilityTimeout(120))
}

// Creates a filter depending on whether the execution is local (dummy filter) or not (Pyth filter)
func newFilterFunc(cfg *config.ServiceConfiguration) queue.FilterConsumeFunc {
	if cfg.P2pNetwork == config.P2pMainNet {
		return queue.PythFilter
	}
	return queue.NonFilter
}

// Creates a metrics depending on whether the execution is local (dummy metrics) or not (Prometheus metrics)
func newMetrics(cfg *config.ServiceConfiguration) metrics.Metrics {
	if !cfg.MetricsEnabled {
		return metrics.NewDummyMetrics()
	}
	return metrics.NewPrometheusMetrics(cfg.Environment)
}

func newAlertClient(cfg *config.ServiceConfiguration) (alert.AlertClient, error) {
	if !cfg.AlertEnabled {
		return alert.NewDummyClient(), nil
	}

	alertConfig := alert.AlertConfig{
		Environment: cfg.Environment,
		ApiKey:      cfg.AlertApiKey,
		Enabled:     cfg.AlertEnabled,
	}

	return alert.NewAlertService(alertConfig, parserAlert.LoadAlerts)
}
