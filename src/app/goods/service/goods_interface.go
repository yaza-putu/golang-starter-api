package service

import (
	"context"
	"github.com/yaza-putu/golang-starter-api/src/app/goods/validation"
	"github.com/yaza-putu/golang-starter-api/src/http/response"
)

type GoodInterface interface {
	Create(ctx context.Context, gds validation.Goods) response.DataApi
	Update(ctx context.Context, id string, gds validation.Goods) response.DataApi
	Delete(ctx context.Context, id string) response.DataApi
	FindById(ctx context.Context, id string) response.DataApi
	Stock(ctx context.Context, id string, n int) response.DataApi
	All(ctx context.Context, page int, take int) response.DataApi
}
