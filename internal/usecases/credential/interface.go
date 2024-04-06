package credential

import (
	"brick/internal/entities"
	"brick/internal/pkg/serror"
	"context"
	"database/sql"
)

type CredentialRepository interface {
	CreateCredential(ctx context.Context, tx *sql.Tx, entityCredential *entities.Credential) *serror.Error
	GetCredentialByClientID(ctx context.Context, entityCredential *entities.Credential) (*entities.Credential, *serror.Error)
	GetCredentialByUserID(ctx context.Context, entityCredential *entities.Credential) (*entities.Credential, *serror.Error)
}
