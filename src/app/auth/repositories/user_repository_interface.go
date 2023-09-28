package repositories

import "yaza/src/app/auth/entities"

type UserInterface interface {
	Create(user entities.User) (entities.User, error)
	Update(id string, user entities.User) error
	Delete(id string) error
	FindOrFail(id string) (entities.User, error)
}
