package presenters

import (
	"brick/internal/adapters/dto/response"
	"brick/internal/entities"
	"brick/internal/usecases/user"

	"github.com/sirupsen/logrus"
)

type userPresenter struct {
	Log *logrus.Logger
}

func NewUserPresenter(log *logrus.Logger) user.UserPresenter {
	return &userPresenter{
		Log: log,
	}
}

func (p *userPresenter) PresentLoginUser(accessToken string) *response.LoginUser {
	return &response.LoginUser{
		AccessToken: accessToken,
	}
}

func (p *userPresenter) PresentCreateUser(entityUser *entities.User) *response.CreateUser {
	return &response.CreateUser{
		ID:    entityUser.ID,
		Email: entityUser.Email,
	}
}

func (p *userPresenter) PresentGetUserInfoByID(entityUser *entities.User, entityCredential *entities.Credential) *response.GetUserInfo {
	return &response.GetUserInfo{
		ID:      entityUser.ID,
		Email:   entityUser.Email,
		Balance: entityUser.Balance,
		Credential: response.Credential{
			ID:           entityCredential.ID,
			ClientID:     entityCredential.ClientID,
			ClientSecret: entityCredential.ClientSecret,
		},
	}
}

func (p *userPresenter) PresentGenerateToken(apiToken string) *response.GenerateToken {
	return &response.GenerateToken{
		ApiToken: apiToken,
	}
}
