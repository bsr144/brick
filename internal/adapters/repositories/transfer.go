package repositories

import (
	"brick/internal/entities"
	"brick/internal/pkg/constvar"
	"brick/internal/pkg/query"
	"brick/internal/pkg/serror"
	"brick/internal/usecases/transfer"
	"context"
	"database/sql"
	"net/http"

	"github.com/sirupsen/logrus"
)

type transferRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewTransferRepository(db *sql.DB, log *logrus.Logger) transfer.TransferRepository {
	return &transferRepository{
		DB:  db,
		Log: log,
	}
}

func (r *transferRepository) CreateTransfer(ctx context.Context, entityTransfer *entities.Transfer) (*entities.Transfer, *serror.Error) {
	err := r.DB.
		QueryRowContext(ctx, query.CREATE_NEW_TRANSFER_WITH_RETURNING_QUERY, entityTransfer.RecipientAccountID, entityTransfer.SenderAccountID, entityTransfer.Amount, entityTransfer.Status).
		Scan(
			&entityTransfer.ID,
			&entityTransfer.RecipientAccountID,
			&entityTransfer.SenderAccountID,
			&entityTransfer.Amount,
			&entityTransfer.Status,
			&entityTransfer.CreatedAt,
		)

	if err != nil {
		r.Log.Errorf("[Transfer][Repository] while tx.QueryRowContext: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return entityTransfer, nil
}

func (r *transferRepository) UpdateTransferStatus(ctx context.Context, entityTransfer *entities.Transfer) *serror.Error {
	_, err := r.DB.
		ExecContext(ctx, query.UPDATE_TRANSFER_STATUS_BY_ID, entityTransfer.ID, entityTransfer.Status)

	if err != nil {
		r.Log.Errorf("[Transfer][Repository] while tx.ExecContext: %s", err.Error())
		return serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return nil
}
