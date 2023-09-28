package services

import (
	"github.com/golang-jwt/jwt"
	"yaza/src/app/auth/entities"
)

type TokenInterface interface {
	Create(user entities.User) (string, string, error)
	Refresh(rToken string) (string, error)
	generateToken(claims jwt.Claims) (string, error)
	generateRefreshToken(claims jwt.Claims) (string, error)
}
