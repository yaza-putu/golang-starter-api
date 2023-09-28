package repositories

import (
	"yaza/src/app/auth/entities"
)

type userRepository struct {
	entity entities.User
}

func NewUserRepository() UserInterface {
	return &userRepository{
		entity: entities.User{},
	}
}

func (u *userRepository) Create(user entities.User) (entities.User, error) {
	return entities.User{}, nil
}

func (u *userRepository) Update(id string, user entities.User) error {
	return nil
}

func (u *userRepository) Delete(id string) error {
	return nil
}

func (u *userRepository) FindOrFail(id string) (entities.User, error) {
	return entities.User{}, nil
}
