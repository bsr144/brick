package repositories

import (
	"brick/internal/entities"
	"brick/internal/pkg/constvar"
	"brick/internal/pkg/query"
	"brick/internal/pkg/serror"
	"brick/internal/usecases/credential"
	"context"
	"database/sql"
	"net/http"

	"github.com/sirupsen/logrus"
)

type credentialRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewCredentialRepository(db *sql.DB, log *logrus.Logger) credential.CredentialRepository {
	return &credentialRepository{
		DB:  db,
		Log: log,
	}
}

func (r *credentialRepository) CreateCredential(ctx context.Context, tx *sql.Tx, entityCredential *entities.Credential) *serror.Error {
	_, err := tx.ExecContext(ctx, query.CREATE_NEW_CREDENTIAL, entityCredential.ClientID, entityCredential.ClientSecret, entityCredential.UserID)
	if err != nil {
		r.Log.Errorf("[Credential][Repository] while tx.ExecContext: %s", err.Error())
		return serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}
	return nil
}

func (r *credentialRepository) GetCredentialByUserID(ctx context.Context, entityCredential *entities.Credential) (*entities.Credential, *serror.Error) {
	err := r.DB.
		QueryRowContext(ctx, query.GET_CREDENTIAL_BY_USER_ID, entityCredential.UserID).
		Scan(&entityCredential.ID,
			&entityCredential.ClientID,
			&entityCredential.ClientSecret,
			&entityCredential.UserID,
		)

	if err == sql.ErrNoRows {
		return entityCredential, nil
	}
	if err != nil {
		r.Log.Errorf("[Credential][Repository] while DB.QueryRowCOntext: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return entityCredential, nil
}

func (r *credentialRepository) GetCredentialByClientID(ctx context.Context, entityCredential *entities.Credential) (*entities.Credential, *serror.Error) {
	err := r.DB.
		QueryRowContext(ctx, query.GET_CREDENTIAL_BY_CLIENT_ID, entityCredential.ClientID).
		Scan(&entityCredential.ID,
			&entityCredential.ClientID,
			&entityCredential.ClientSecret,
			&entityCredential.UserID,
		)

	if err == sql.ErrNoRows {
		return entityCredential, nil
	}

	if err != nil {
		r.Log.Errorf("[Credential][Repository] while DB.QueryRowContext: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return entityCredential, nil
}
