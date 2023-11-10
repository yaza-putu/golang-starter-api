package service

import (
	"context"
	"github.com/yaza-putu/golang-starter-api/src/http/response"
)

type AuthInterface interface {
	Login(ctx context.Context, email string, password string) response.DataApi
	Refresh(ctx context.Context, oToken string) response.DataApi
}
