package services

import (
	"context"
	"yaza/src/app/auth/repositories"
	"yaza/src/http/response"
	"yaza/src/utils"
)

type authService struct {
	tokenService   TokenInterface
	userRepository repositories.UserInterface
}

func NewAuthService() AuthInterface {
	return &authService{
		userRepository: repositories.NewUserRepository(),
		tokenService:   NewToken(),
	}
}

func (a *authService) Login(ctx context.Context, email string, password string) response.DataApi {
	rc := make(chan response.DataApi)
	go func() {

		dUser, err := a.userRepository.FindByEmail(ctx, email)
		if err != nil {
			rc <- response.Api(response.SetCode(500), response.SetMessage(err))
		}

		if utils.BcryptCheck(password, dUser.Password) {
			// generate token
			token, refresh, err := a.tokenService.Create(ctx, dUser)
			if err != nil {
				rc <- response.Api(response.SetCode(500), response.SetMessage(err))
			}

			rc <- response.Api(response.SetCode(200), response.SetMessage("Generate token successfully"), response.SetData(map[string]string{
				"access_token":  token,
				"refresh_token": refresh,
			}))
		} else {
			rc <- response.Api(response.SetCode(401), response.SetMessage("Credential not authorized"))
		}
		close(rc)
	}()

	for {
		select {
		case <-ctx.Done():
			return response.Api(response.SetCode(408), response.SetMessage("Request timeout or canceled by user"))
		case res := <-rc:
			return response.Api(response.SetCode(res.Code), response.SetData(res.Data))
		}
	}
}

func (a *authService) Refresh(ctx context.Context, oToken string) response.DataApi {
	rc := make(chan response.DataApi)

	go func() {
		token, err := a.tokenService.Refresh(ctx, oToken)
		if err != nil {
			rc <- response.Api(response.SetCode(500), response.SetMessage(err))
		}

		rc <- response.Api(response.SetCode(200), response.SetData(token))
	}()

	for {
		select {
		case <-ctx.Done():
			return response.Api(response.SetCode(408), response.SetMessage("Request timeout or canceled by user"))
		case res := <-rc:
			return response.Api(response.SetCode(res.Code), response.SetData(res.Data))
		}
	}
}
