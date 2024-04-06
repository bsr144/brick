package user

import (
	"brick/internal/adapters/dto/request"
	"brick/internal/adapters/dto/response"
	"brick/internal/entities"
	"brick/internal/pkg/serror"
	"context"
	"database/sql"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx *sql.Tx, entityUser *entities.User) (*entities.User, *serror.Error)
	GetUserByEmail(ctx context.Context, entityUser *entities.User) (*entities.User, *serror.Error)
	GetUserByID(ctx context.Context, entityUser *entities.User) (*entities.User, *serror.Error)
}

type UserPresenter interface {
	PresentCreateUser(entityUser *entities.User) *response.CreateUser
	PresentLoginUser(token string) *response.LoginUser
	PresentGetUserInfoByID(entityUser *entities.User, entityCredential *entities.Credential) *response.GetUserInfo
	PresentGenerateToken(token string) *response.GenerateToken
}

type UserUseCase interface {
	CreateUser(ctx context.Context, createUserRequest *request.CreateUser) (*response.CreateUser, *serror.Error)
	LoginUser(ctx context.Context, loginUserRequest *request.LoginUser) (*response.LoginUser, *serror.Error)
	GetUserByID(ctx context.Context, getInfoRequest *request.GetUserInfo) (*response.GetUserInfo, *serror.Error)
	GenerateToken(ctx context.Context, generateTokenRequest *request.GenerateToken) (*response.GenerateToken, *serror.Error)
}
