package repositories

import (
	"brick/internal/entities"
	"brick/internal/pkg/constvar"
	"brick/internal/pkg/query"
	"brick/internal/pkg/serror"
	"brick/internal/usecases/user"
	"context"
	"database/sql"
	"net/http"

	"github.com/sirupsen/logrus"
)

type userRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewUserRepository(db *sql.DB, log *logrus.Logger) user.UserRepository {
	return &userRepository{
		DB:  db,
		Log: log,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, tx *sql.Tx, entityUser *entities.User) (*entities.User, *serror.Error) {
	err := tx.
		QueryRowContext(ctx, query.CREATE_NEW_USER_WITH_RETURNING_QUERY, entityUser.Email, entityUser.Password, entityUser.Salt, entityUser.Balance).
		Scan(
			&entityUser.ID,
			&entityUser.Email,
		)

	if err != nil {
		r.Log.Errorf("[User][Repository] while tx.QueryRowContext: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return entityUser, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, entityUser *entities.User) (*entities.User, *serror.Error) {
	err := r.DB.
		QueryRowContext(ctx, query.GET_USER_BY_EMAIL_QUERY, entityUser.Email).
		Scan(&entityUser.ID,
			&entityUser.Email,
			&entityUser.Password,
			&entityUser.Balance,
			&entityUser.Salt,
			&entityUser.CreatedAt,
			&entityUser.UpdatedAt,
			&entityUser.DeletedAt,
		)

	if err == sql.ErrNoRows {
		return entityUser, nil
	}

	if err != nil {
		r.Log.Errorf("[User][Repository] while DB.QueryRowContext: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return entityUser, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, entityUser *entities.User) (*entities.User, *serror.Error) {
	err := r.DB.
		QueryRowContext(ctx, query.GET_USER_BY_ID, entityUser.ID).
		Scan(&entityUser.ID,
			&entityUser.Email,
			&entityUser.Balance,
		)

	if err == sql.ErrNoRows {
		return entityUser, nil
	}

	if err != nil {
		r.Log.Errorf("[User][Repository] while DB.QueryRowContext: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return entityUser, nil
}
