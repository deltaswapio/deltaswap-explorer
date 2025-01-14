package main

import (
	"context"
	"log"
	"os"

	"github.com/deltaswapio/deltaswap-explorer/common/dbutil"
	"github.com/deltaswapio/deltaswap-explorer/common/logger"
	"github.com/deltaswapio/deltaswap-explorer/common/prices"
	"github.com/deltaswapio/deltaswap-explorer/jobs/config"
	"github.com/deltaswapio/deltaswap-explorer/jobs/internal/coingecko"
	"github.com/deltaswapio/deltaswap-explorer/jobs/jobs"
	"github.com/deltaswapio/deltaswap-explorer/jobs/jobs/notional"
	"github.com/deltaswapio/deltaswap-explorer/jobs/jobs/report"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type exitCode int

func main() {
	defer handleExit()
	context := context.Background()

	// get the config
	cfg, errConf := config.New(context)
	if errConf != nil {
		log.Fatal("error creating config", errConf)
	}

	logger := logger.New("deltaswap-explorer-jobs", logger.WithLevel(cfg.LogLevel))
	logger.Info("started job execution", zap.String("job_id", cfg.JobID))

	var err error
	switch cfg.JobID {
	case jobs.JobIDNotional:
		nCfg, errCfg := config.NewNotionalConfiguration(context)
		if errCfg != nil {
			log.Fatal("error creating config", errCfg)
		}
		notionalJob := initNotionalJob(context, nCfg, logger)
		err = notionalJob.Run()
	case jobs.JobIDTransferReport:
		aCfg, errCfg := config.NewTransferReportConfiguration(context)
		if errCfg != nil {
			log.Fatal("error creating config", errCfg)
		}
		transferReport := initTransferReportJob(context, aCfg, logger)
		err = transferReport.Run(context)

	default:
		logger.Fatal("Invalid job id", zap.String("job_id", cfg.JobID))
	}

	if err != nil {
		logger.Error("failed job execution", zap.String("job_id", cfg.JobID), zap.Error(err))
	} else {
		logger.Info("finish job execution successfully", zap.String("job_id", cfg.JobID))
	}

}

// initNotionalJob initializes notional job.
func initNotionalJob(ctx context.Context, cfg *config.NotionalConfiguration, logger *zap.Logger) *notional.NotionalJob {
	// init coingecko api client.
	api := coingecko.NewCoingeckoAPI(cfg.CoingeckoURL)
	// init redis client.
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.CacheURL})
	// create notional job.
	notionalJob := notional.NewNotionalJob(api, redisClient, cfg.CachePrefix, cfg.P2pNetwork, cfg.NotionalChannel, logger)
	return notionalJob
}

// initTransferReportJob initializes transfer report job.
func initTransferReportJob(ctx context.Context, cfg *config.TransferReportConfiguration, logger *zap.Logger) *report.TransferReportJob {
	//setup DB connection
	db, err := dbutil.Connect(ctx, logger, cfg.MongoURI, cfg.MongoDatabase, false)
	if err != nil {
		logger.Fatal("Failed to connect MongoDB", zap.Error(err))
	}
	pricesCache := prices.NewCoinPricesCache(cfg.PricesPath)
	pricesCache.InitCache()
	return report.NewTransferReportJob(db.Database, cfg.PageSize, pricesCache, cfg.OutputPath, logger)
}

func handleExit() {
	if r := recover(); r != nil {
		if e, ok := r.(exitCode); ok {
			os.Exit(int(e))
		}
		panic(r) // not an Exit, bubble up
	}
}
