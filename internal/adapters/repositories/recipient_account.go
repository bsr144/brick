package repositories

import (
	"brick/internal/entities"
	"brick/internal/pkg/constvar"
	"brick/internal/pkg/query"
	"brick/internal/pkg/serror"
	"brick/internal/usecases/recipientAccount"
	"context"
	"database/sql"
	"net/http"

	"github.com/sirupsen/logrus"
)

type recipientAccountRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewRecipientAccountRepository(db *sql.DB, log *logrus.Logger) recipientAccount.RecipientAccountRepository {
	return &recipientAccountRepository{
		DB:  db,
		Log: log,
	}
}

func (r *recipientAccountRepository) CreateRecipientAccount(ctx context.Context, entityRecipientAccount *entities.RecipientAccount) (*entities.RecipientAccount, *serror.Error) {
	err := r.DB.
		QueryRowContext(
			ctx,
			query.CREATE_NEW_RECIPIENT_ACCOUNT_WITH_RETURNING_QUERY,
			entityRecipientAccount.AccountNumber,
			entityRecipientAccount.AccountName,
			entityRecipientAccount.BankCode,
			entityRecipientAccount.BankName,
			entityRecipientAccount.VerificationStatus,
		).
		Scan(
			&entityRecipientAccount.ID,
			&entityRecipientAccount.AccountNumber,
			&entityRecipientAccount.AccountName,
			&entityRecipientAccount.BankCode,
			&entityRecipientAccount.BankName,
			&entityRecipientAccount.VerificationStatus,
		)

	if err != nil {
		r.Log.Errorf("[RecipientAccount][Repository] while DB.QueryRowContext: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return entityRecipientAccount, nil
}

func (r *recipientAccountRepository) UpdateRecipientAccountByID(ctx context.Context, entityRecipientAccount *entities.RecipientAccount) (*entities.RecipientAccount, *serror.Error) {
	err := r.DB.
		QueryRowContext(
			ctx,
			query.UPDATE_RECIPIENT_ACCOUNT_BY_ID,
			entityRecipientAccount.ID,
			entityRecipientAccount.AccountNumber,
			entityRecipientAccount.AccountName,
			entityRecipientAccount.BankCode,
			entityRecipientAccount.BankName,
			entityRecipientAccount.VerificationStatus).
		Scan(&entityRecipientAccount.ID,
			&entityRecipientAccount.AccountNumber,
			&entityRecipientAccount.AccountName,
			&entityRecipientAccount.BankCode,
			&entityRecipientAccount.BankName,
			&entityRecipientAccount.VerificationStatus,
			&entityRecipientAccount.LastVerifiedAt)
	if err == sql.ErrNoRows {
		return entityRecipientAccount, nil
	}

	if err != nil {
		r.Log.Errorf("[RecipientAccount][Repository] while DB.QueryRowContext: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return entityRecipientAccount, nil
}

func (r *recipientAccountRepository) GetRecipientAccountByBankCodeAndAccountNumber(ctx context.Context, entityRecipientAccount *entities.RecipientAccount) (*entities.RecipientAccount, *serror.Error) {
	err := r.DB.
		QueryRowContext(ctx, query.GET_RECIPIENT_ACCOUNT_BY_BANK_CODE_AND_BANK_ACCOUNT_NUMBER, entityRecipientAccount.BankCode, entityRecipientAccount.AccountNumber).
		Scan(
			&entityRecipientAccount.ID,
			&entityRecipientAccount.AccountNumber,
			&entityRecipientAccount.AccountName,
			&entityRecipientAccount.BankCode,
			&entityRecipientAccount.BankName,
			&entityRecipientAccount.VerificationStatus,
			&entityRecipientAccount.LastVerifiedAt,
		)
	if err == sql.ErrNoRows {
		return entityRecipientAccount, nil
	}

	if err != nil {
		r.Log.Errorf("[RecipientAccount][Repository] while DB.QueryRowContext: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return entityRecipientAccount, nil
}
