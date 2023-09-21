package services

import (
	"github.com/golang-jwt/jwt"
	"yaza/src/app/auth/model"
)

type TokenInterface interface {
	Create(user model.User) (string, string, error)
	Refresh(rToken string) (string, error)
	generateToken(claims jwt.Claims) (string, error)
	generateRefreshToken(claims jwt.Claims) (string, error)
}
