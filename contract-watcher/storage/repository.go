package storage

import (
	"context"
	"errors"
	"time"

	"github.com/deltaswapio/deltaswap-explorer/common/client/alert"
	cwAlert "github.com/deltaswapio/deltaswap-explorer/contract-watcher/internal/alert"
	"github.com/deltaswapio/deltaswap-explorer/contract-watcher/internal/metrics"
	sdk "github.com/deltaswapio/deltaswap/sdk/vaa"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// repository errors
var ErrDocNotFound = errors.New("NOT FOUND")

// Repository definitions.
type Repository struct {
	db          *mongo.Database
	log         *zap.Logger
	metrics     metrics.Metrics
	alerts      alert.AlertClient
	collections struct {
		watcherBlock       *mongo.Collection
		globalTransactions *mongo.Collection
	}
}

// NewRepository create a new respository instance.
func NewRepository(db *mongo.Database, metrics metrics.Metrics, alerts alert.AlertClient, log *zap.Logger) *Repository {
	return &Repository{db, log, metrics, alerts, struct {
		watcherBlock       *mongo.Collection
		globalTransactions *mongo.Collection
	}{
		watcherBlock:       db.Collection("watcherBlock"),
		globalTransactions: db.Collection("globalTransactions"),
	}}
}

func indexedAt(t time.Time) IndexingTimestamps {
	return IndexingTimestamps{
		IndexedAt: t,
	}
}

func (s *Repository) UpsertGlobalTransaction(ctx context.Context, chainID sdk.ChainID, globalTx TransactionUpdate) error {
	update := bson.M{
		"$set":         globalTx,
		"$setOnInsert": indexedAt(time.Now()),
		"$inc":         bson.D{{Key: "revision", Value: 1}},
	}

	_, err := s.collections.globalTransactions.UpdateByID(ctx, globalTx.ID, update, options.Update().SetUpsert(true))
	if err != nil {
		s.log.Error("Error inserting global transaction", zap.Error(err))
		// send alert when exists an error saving ptth vaa.
		alertContext := alert.AlertContext{
			Details: globalTx.ToMap(),
			Error:   err,
		}
		s.alerts.CreateAndSend(ctx, cwAlert.ErrorSaveDestinationTx, alertContext)
		return err
	}
	s.metrics.IncDestinationTrxSaved(chainID)

	return err

}

func (s *Repository) GetGlobalTransactionByID(ctx context.Context, id string) (TransactionUpdate, error) {
	var tx TransactionUpdate
	err := s.collections.globalTransactions.FindOne(ctx, bson.M{"_id": id}).Decode(&tx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return tx, ErrDocNotFound
		}
		return tx, err
	}
	return tx, nil
}

func (s *Repository) UpdateWatcherBlock(ctx context.Context, chainID sdk.ChainID, watcherBlock WatcherBlock) error {
	update := bson.M{
		"$set":         watcherBlock,
		"$setOnInsert": indexedAt(time.Now()),
	}
	s.metrics.SetCurrentBlock(chainID, uint64(watcherBlock.BlockNumber))
	_, err := s.collections.watcherBlock.UpdateByID(ctx, watcherBlock.ID, update, options.Update().SetUpsert(true))
	if err != nil {
		s.log.Error("Error inserting watcher block", zap.Error(err))
		return err
	}
	return err
}

func (s *Repository) GetCurrentBlock(ctx context.Context, blockchain string, defaultBlock int64) (int64, error) {
	var block WatcherBlock
	err := s.collections.watcherBlock.FindOne(ctx, bson.M{"_id": blockchain}).Decode(&block)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return defaultBlock, nil
		}
		return 0, err
	}
	return block.BlockNumber, nil
}
