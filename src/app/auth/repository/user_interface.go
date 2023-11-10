package repository

import (
	"context"
	"github.com/yaza-putu/golang-starter-api/src/app/auth/entity"
)

type UserInterface interface {
	FindByEmail(ctx context.Context, email string) (entity.User, error)
}
