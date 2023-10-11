package services

import (
	"context"
	"yaza/src/http/response"
)

type AuthInterface interface {
	Login(ctx context.Context, email string, password string) response.DataApi
	Refresh(ctx context.Context, oToken string) response.DataApi
}
