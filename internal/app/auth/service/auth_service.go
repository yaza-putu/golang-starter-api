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
	Login(ctx context.Context, email string, password string) response.DataApi
	Refresh(ctx context.Context, oToken string) response.DataApi
}

type authService struct {
	tokenService   Token
	userRepository repository.User
}

func NewAuth(u repository.User, t Token) Auth {
	return &authService{
		userRepository: u,
		tokenService:   t,
	}
}

func (a *authService) Login(ctx context.Context, email string, password string) response.DataApi {
	rc := make(chan response.DataApi)
	go func() {

		dUser, err := a.userRepository.FindByEmail(ctx, email)
		if err != nil {
			if err.Error() == "record not found" {
				rc <- response.Api(response.SetCode(401), response.SetMessage("Kredensial tidak valid"))
			}
			rc <- response.Api(response.SetCode(500), response.SetError(err))
		}

		if encrypt.BcryptCheck(password, dUser.Password) {
			// generate token
			token, refresh, err := a.tokenService.Create(ctx, dUser)
			if err != nil {
				rc <- response.Api(response.SetCode(500), response.SetError(err))
			}

			rc <- response.Api(response.SetCode(200), response.SetMessage("Generate token successfully"), response.SetData(map[string]string{
				"access_token":  token,
				"refresh_token": refresh,
			}))
		} else {
			rc <- response.Api(response.SetCode(401), response.SetMessage("Kredensial tidak valid"))
		}
		close(rc)
	}()

	for {
		select {
		case <-ctx.Done():
			return response.TimeOut()
		case res := <-rc:
			return res
		}
	}
}

func (a *authService) Refresh(ctx context.Context, oToken string) response.DataApi {
	rc := make(chan response.DataApi)

	go func() {
		token, err := a.tokenService.Refresh(ctx, oToken)
		if err != nil {
			if err.Error() == "Token is expired" {
				rc <- response.Api(response.SetCode(401), response.SetMessage("Token is expired"), response.SetError(err))
			} else {
				logger.New(err, logger.SetType(logger.ERROR))
				rc <- response.Api(response.SetCode(500), response.SetError(err))
			}
		}

		rc <- response.Api(response.SetCode(200), response.SetData(token))
	}()

	for {
		select {
		case <-ctx.Done():
			return response.TimeOut()
		case res := <-rc:
			return res
		}
	}
}
