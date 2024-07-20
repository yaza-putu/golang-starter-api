package service

import (
	"context"

	"github.com/yaza-putu/golang-starter-api/internal/app/auth/repository"
	"github.com/yaza-putu/golang-starter-api/internal/http/response"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/encrypt"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
)

// Auth / **************************************************************
type Auth interface {
	Login(email string, password string) response.DataApi
	Refresh(oToken string) response.DataApi
}

type authService struct {
	tokenService   Token
	userRepository repository.User
}

func NewAuth(u repository.User, t Token) *authService {
	return &authService{
		userRepository: u,
		tokenService:   t,
	}
}

func (a *authService) Login(email string, password string) response.DataApi {
	dUser, err := a.userRepository.FindByEmail(context.Background(), email)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Api(response.SetCode(401), response.SetMessage("Invalid credentials"))
		}
		return response.Api(response.SetCode(500), response.SetError(err))
	}

	if encrypt.BcryptCheck(password, dUser.Password) {
		// generate token
		token, refresh, err := a.tokenService.Create(dUser)
		if err != nil {
			return response.Api(response.SetCode(500), response.SetError(err))
		}

		return response.Api(response.SetCode(200), response.SetMessage("Generate token successfully"), response.SetData(map[string]string{
			"access_token":  token,
			"refresh_token": refresh,
		}))
	} else {
		return response.Api(response.SetCode(401), response.SetMessage("Invalid credentials"))
	}
}

func (a *authService) Refresh(oToken string) response.DataApi {
	token, err := a.tokenService.Refresh(oToken)
	if err != nil {
		if err.Error() == "Token is expired" {
			return response.Api(response.SetCode(401), response.SetMessage("Token is expired"), response.SetError(err))
		} else {
			logger.New(err, logger.SetType(logger.ERROR))
			return response.Api(response.SetCode(500), response.SetError(err))
		}
	}

	return response.Api(response.SetCode(200), response.SetData(token))
}
