package repository

import (
	"context"
	"github.com/yaza-putu/golang-starter-api/src/app/category/entity"
	"github.com/yaza-putu/golang-starter-api/src/database"
	redis_client "github.com/yaza-putu/golang-starter-api/src/redis"
	"github.com/yaza-putu/golang-starter-api/src/utils"
	"strings"
)

type categoryRepository struct {
	entity   entity.Category
	entities entity.Categories
}

func NewCategory() CategoryInterface {
	return &categoryRepository{
		entity:   entity.Category{},
		entities: entity.Categories{},
	}
}

func (c *categoryRepository) Create(ctx context.Context, cat entity.Category) (entity.Category, error) {
	cat.ID = utils.Uid(13)
	cat.Name = strings.ToTitle(cat.Name)
	r := database.Instance.WithContext(ctx).Create(&cat)
	redis_client.Set(ctx, cat.ID, cat)
	return cat, r.Error
}

func (c *categoryRepository) Update(ctx context.Context, id string, cat entity.Category) error {
	cat.Name = strings.ToTitle(cat.Name)

	d := cat
	d.ID = id
	redis_client.Set(ctx, id, d)
	return database.Instance.WithContext(ctx).Where("id = ?", id).Updates(&cat).Error
}

func (c *categoryRepository) Delete(ctx context.Context, id string) error {
	redis_client.Del(context.Background(), id)
	return database.Instance.WithContext(ctx).Where("id = ?", id).Delete(&c.entity).Error
}

func (c *categoryRepository) FindById(ctx context.Context, id string) (entity.Category, error) {
	e := c.entity
	r := database.Instance.WithContext(ctx).Where("id = ?", id).First(&e)
	if r.Error == nil {
		redis_client.FindSet(ctx, id, e)
	}
	return e, r.Error
}

func (c *categoryRepository) All(ctx context.Context, page int, take int) (utils.Pagination, error) {
	e := c.entities
	var pagination utils.Pagination
	var totalRow int64

	r := database.Instance.WithContext(ctx).Model(&e)
	r.Count(&totalRow)
	r.Scopes(pagination.Paginate(page, take))
	r.Find(&e)

	pagination.Rows = e
	pagination.CalculatePage(float64(totalRow))

	return pagination, r.Error
}
