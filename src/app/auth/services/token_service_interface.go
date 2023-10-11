package services

import (
	"context"
	"github.com/golang-jwt/jwt"
	"yaza/src/app/auth/entities"
)

type TokenInterface interface {
	Create(ctx context.Context, user entities.User) (string, string, error)
	Refresh(ctx context.Context, rToken string) (string, error)
	generateToken(claims jwt.Claims) (string, error)
	generateRefreshToken(claims jwt.Claims) (string, error)
}
