package repositories

import (
	"context"
	"yaza/src/app/auth/entities"
)

type UserInterface interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	Update(ctx context.Context, id string, user entities.User) error
	Delete(ctx context.Context, id string) error
	FindOrFail(ctx context.Context, id string) (entities.User, error)
	FindByEmail(ctx context.Context, email string) (entities.User, error)
}
