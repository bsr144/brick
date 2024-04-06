package user

import (
	"brick/internal/adapters/dto/request"
	"brick/internal/adapters/dto/response"
	"brick/internal/entities"
	"brick/internal/pkg/constvar"
	"brick/internal/pkg/protection"
	"brick/internal/pkg/serror"
	"brick/internal/usecases/credential"
	"context"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type usecase struct {
	UserRepository       UserRepository
	UserPresenter        UserPresenter
	CredentialRepository credential.CredentialRepository
	DB                   *sql.DB
	Log                  *logrus.Logger
}

func NewUserUsecase(userRepository UserRepository, credentialRepository credential.CredentialRepository, presenter UserPresenter, db *sql.DB, log *logrus.Logger) UserUseCase {
	return &usecase{
		UserRepository:       userRepository,
		CredentialRepository: credentialRepository,
		UserPresenter:        presenter,
		DB:                   db,
		Log:                  log,
	}
}

func (u *usecase) CreateUser(ctx context.Context, createUserRequest *request.CreateUser) (*response.CreateUser, *serror.Error) {
	user := new(entities.User)
	user.Email = createUserRequest.Email
	user, serr := u.UserRepository.GetUserByEmail(ctx, user)
	if serr != nil {
		u.Log.Errorf("[User][Usecase] while UserRepository.GetUserByEmail: %s", serr.Error())
		return nil, serr
	}
	if user.IsExist() {
		u.Log.Errorf("[User][Usecase] while checking user: %s", constvar.USER_IS_EXIST_ERROR)
		return nil, serror.NewError(http.StatusBadRequest, 0, constvar.USER_IS_EXIST_ERROR, constvar.USER_IS_EXIST_ERROR)
	}

	user.Password = createUserRequest.Password
	user.Balance = createUserRequest.Balance
	user.Salt, serr = protection.GenerateSalt()
	if serr != nil {
		u.Log.Errorf("[User][Usecase] while protection.GenerateSalt: %s", serr.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, serr.Error())
	}
	user.Password = protection.HashPassword(user.Password, user.Salt)

	tx, err := u.DB.Begin()
	if err != nil {
		u.Log.Errorf("[User][Usecase] while DB.Begin: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, serr.Error())
	}

	defer tx.Rollback()

	user, serr = u.UserRepository.CreateUser(ctx, tx, user)
	if serr != nil {
		u.Log.Errorf("[User][Usecase] while UserRepository.CreateUser: %s", serr.Error())
		return nil, serr
	}

	credential := new(entities.Credential)
	credential.UserID = user.ID
	credential.ClientID = uuid.New().String()
	credential.ClientSecret = uuid.New().String()

	serr = u.CredentialRepository.CreateCredential(ctx, tx, credential)

	err = tx.Commit()
	if err != nil {
		return nil, serror.NewError(http.StatusInternalServerError, 0, constvar.SERVER_INFO_ERROR, serr.Error())
	}

	return u.UserPresenter.PresentCreateUser(user), nil
}

func (u *usecase) LoginUser(ctx context.Context, loginUserRequest *request.LoginUser) (*response.LoginUser, *serror.Error) {
	user := new(entities.User)
	user.Email = loginUserRequest.Email

	user, serr := u.UserRepository.GetUserByEmail(ctx, user)

	if serr != nil {
		u.Log.Errorf("[User][Usecase] while UserRepository.GetUserByEmail: %s", serr.Error())
		return nil, serr
	}
	if user.IsNotExist() {
		u.Log.Errorf("[User][Usecase] while checking user: %s", constvar.USER_NOT_EXIST_ERROR)
		return nil, serror.NewError(http.StatusNotFound, 0, constvar.USER_NOT_EXIST_ERROR, constvar.USER_NOT_EXIST_ERROR)
	}

	isPasswordValid := protection.ComparePassword(user.Password, user.Salt, loginUserRequest.Password)
	if !isPasswordValid {
		u.Log.Errorf("[User][Usecase] while validate password: %s", "password is not valid")
		return nil, serror.NewError(http.StatusBadRequest, 0, "the password given doesn't match", "the password given doesn't match")
	}

	userAccessToken, serr := protection.GenerateToken(user.ID)
	if serr != nil {
		u.Log.Errorf("[User][Usecase] while protection.GenerateToken: %s", serr.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, "error while generating user token", serr.Error())
	}

	return u.UserPresenter.PresentLoginUser(userAccessToken), nil
}

func (u *usecase) GetUserByID(ctx context.Context, getInfoRequest *request.GetUserInfo) (*response.GetUserInfo, *serror.Error) {
	user := new(entities.User)
	user.ID = getInfoRequest.ID

	user, serr := u.UserRepository.GetUserByID(ctx, user)
	if serr != nil {
		u.Log.Errorf("[User][Usecase] while UserRepository.GetUserByID: %s", serr.Error())
		return nil, serr
	}
	if user.IsNotExist() {
		u.Log.Errorf("[User][Usecase] while checking user: %s", constvar.USER_NOT_EXIST_ERROR)
		return nil, serror.NewError(http.StatusNotFound, 0, constvar.USER_NOT_EXIST_ERROR, constvar.USER_NOT_EXIST_ERROR)
	}

	credential := new(entities.Credential)
	credential.UserID = user.ID
	credential, serr = u.CredentialRepository.GetCredentialByUserID(ctx, credential)
	if serr != nil {
		u.Log.Errorf("[User][Usecase] while CredentialRepository.GetCredentialByUserID: %s", serr.Error())
		return nil, serr
	}

	return u.UserPresenter.PresentGetUserInfoByID(user, credential), nil
}

func (u *usecase) GenerateToken(ctx context.Context, generateTokenRequest *request.GenerateToken) (*response.GenerateToken, *serror.Error) {
	credential := new(entities.Credential)
	credential.ClientID = generateTokenRequest.ClientID

	credential, serr := u.CredentialRepository.GetCredentialByClientID(ctx, credential)
	if serr != nil {
		u.Log.Errorf("[User][Usecase] while CredentialRepository.GetCredentialByClientID")
		return nil, serr
	}

	if credential.IsNotExist() {
		u.Log.Errorf("[User][Usecase] while checking credential: %s", "client id not found")
		return nil, serror.NewError(http.StatusUnauthorized, 0, constvar.UNAUTHORIZED_ERROR, "client id is not valid")
	}

	if generateTokenRequest.ClientSecret != credential.ClientSecret {
		u.Log.Errorf("[User][Usecase] while comparing requested client secret with credential client secret")
		return nil, serror.NewError(http.StatusUnauthorized, 0, "Unauthorized", "Unauthorized")
	}

	apiToken, serr := protection.GenerateApiToken(credential.ClientID, credential.UserID)
	if serr != nil {
		u.Log.Errorf("[User][Usecase] while protection.GenerateApiToken: %s", serr.Error())
		return nil, serror.NewError(http.StatusInternalServerError, 0, "error while generating api token", serr.Error())
	}

	return u.UserPresenter.PresentGenerateToken(apiToken), nil
}
