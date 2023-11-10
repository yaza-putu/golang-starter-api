package repository

import (
	"context"
	"github.com/yaza-putu/golang-starter-api/src/app/category/entity"
	"github.com/yaza-putu/golang-starter-api/src/utils"
)

type CategoryInterface interface {
	All(ctx context.Context, page int, take int) (utils.Pagination, error)
	Create(ctx context.Context, cat entity.Category) (entity.Category, error)
	Update(ctx context.Context, id string, user entity.Category) error
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (entity.Category, error)
}
