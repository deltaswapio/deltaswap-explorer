package watcher

import (
	"context"
	"errors"

	"github.com/deltaswapio/deltaswap-explorer/common/domain"
	"github.com/deltaswapio/deltaswap-explorer/contract-watcher/storage"
	sdk "github.com/deltaswapio/deltaswap/sdk/vaa"
	"go.uber.org/zap"
)

var (
	ErrTxfailedCannotBeUpdated = errors.New("tx with status failed can not be updated because exists a confirmed tx for the same vaa ID")
	ErrTxUnknowCannotBeUpdated = errors.New("tx with status unknown can not be updated because exists a tx (confirmed|failed) for the same vaa ID")
	ErrInvalidTxStatus         = errors.New("invalid tx status")
)

type FuncGetGlobalTransactionById func(ctx context.Context, id string) (storage.TransactionUpdate, error)

func updateGlobalTransaction(ctx context.Context, chainID sdk.ChainID, tx storage.TransactionUpdate, r *storage.Repository, log *zap.Logger) {
	updateGlobalTx, err := checkTxShouldBeUpdated(ctx, tx, r.GetGlobalTransactionByID)
	if !updateGlobalTx {
		log.Info("tx can not be updated",
			zap.String("id", tx.ID),
			zap.String("txHash", tx.Destination.TxHash),
			zap.String("status", tx.Destination.Status),
			zap.Error(err))
		return
	}

	err = r.UpsertGlobalTransaction(ctx, chainID, tx)
	if err != nil {
		log.Error("cannot save redeemed tx", zap.Error(err))
	} else {
		log.Info("saved redeemed tx", zap.String("vaa", tx.ID))
	}
}

// checkTxShouldBeUpdated checks if the transaction should be updated.
func checkTxShouldBeUpdated(ctx context.Context, tx storage.TransactionUpdate, getGlobalTransactionByIDFunc FuncGetGlobalTransactionById) (bool, error) {
	switch tx.Destination.Status {
	case domain.DstTxStatusConfirmed:
		return true, nil
	case domain.DstTxStatusFailedToProcess:
		// check if the transaction exists from the same vaa ID.
		oldTx, err := getGlobalTransactionByIDFunc(ctx, tx.ID)
		if err != nil {
			return true, nil
		}
		// if the transaction was already confirmed, then no update it.
		if oldTx.Destination.Status == domain.DstTxStatusConfirmed {
			return false, ErrTxfailedCannotBeUpdated
		}
		return true, nil
	case domain.DstTxStatusUnkonwn:
		// check if the transaction exists from the same vaa ID.
		oldTx, err := getGlobalTransactionByIDFunc(ctx, tx.ID)
		if err != nil {
			return true, nil
		}
		// if the transaction was already confirmed or failed to process, then no update it.
		if oldTx.Destination.Status == domain.DstTxStatusConfirmed || oldTx.Destination.Status == domain.DstTxStatusFailedToProcess {
			return false, ErrTxUnknowCannotBeUpdated
		}
		return true, nil
	default:
		return false, ErrInvalidTxStatus
	}
}
