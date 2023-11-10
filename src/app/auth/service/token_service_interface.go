package service

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/yaza-putu/golang-starter-api/src/app/auth/entity"
)

type TokenInterface interface {
	Create(ctx context.Context, user entity.User) (string, string, error)
	Refresh(ctx context.Context, rToken string) (string, error)
	generateToken(claims jwt.Claims) (string, error)
	generateRefreshToken(claims jwt.Claims) (string, error)
}
