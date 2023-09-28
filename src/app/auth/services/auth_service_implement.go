package services

import (
	"yaza/src/app/auth/repositories"
	"yaza/src/http/response"
)

type authService struct {
	TokenInterface
	userRepository repositories.UserInterface
}

func NewAuthService() AuthInterface {
	return &authService{
		userRepository: repositories.NewUserRepository(),
	}
}

func (a *authService) Login(email string, password string) response.DataApi {
	return response.DataApi{}
}

func (a *authService) Refresh(oToken string) response.DataApi {
	return response.DataApi{}
}
