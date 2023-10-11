package repositories

import (
	"context"
	"yaza/src/app/auth/entities"
	"yaza/src/database"
)

type userRepository struct {
	entity entities.User
}

func NewUserRepository() UserInterface {
	return &userRepository{
		entity: entities.User{},
	}
}

func (u *userRepository) Create(ctx context.Context, user entities.User) (entities.User, error) {
	return entities.User{}, nil
}

func (u *userRepository) Update(ctx context.Context, id string, user entities.User) error {
	return nil
}

func (u *userRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (u *userRepository) FindOrFail(ctx context.Context, id string) (entities.User, error) {
	return entities.User{}, nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	e := u.entity
	r := database.Instance.Where("email = ?", email).First(&e)
	return e, r.Error
}
