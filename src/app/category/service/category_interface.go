package service

import (
	"context"
	"github.com/yaza-putu/golang-starter-api/src/app/category/validation"
	"github.com/yaza-putu/golang-starter-api/src/http/response"
)

type CategoryInterface interface {
	Create(ctx context.Context, cat validation.Category) response.DataApi
	Update(ctx context.Context, id string, cat validation.CategoryWithIgnore) response.DataApi
	Delete(ctx context.Context, id string) response.DataApi
	FindById(ctx context.Context, id string) response.DataApi
	All(ctx context.Context, page int, take int) response.DataApi
}
