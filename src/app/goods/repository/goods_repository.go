package repository

import (
	"context"
	catEntity "github.com/yaza-putu/golang-starter-api/src/app/category/entity"
	"github.com/yaza-putu/golang-starter-api/src/app/goods/entity"
	"github.com/yaza-putu/golang-starter-api/src/database"
	"github.com/yaza-putu/golang-starter-api/src/logger"
	redis_client "github.com/yaza-putu/golang-starter-api/src/redis"
	"github.com/yaza-putu/golang-starter-api/src/utils"
	"gorm.io/gorm"
	"strings"
)

type goodsRepository struct {
	entity   entity.Goods
	entities entity.AllGoods
}

func NewGoods() GoodsInterface {
	return &goodsRepository{
		entity:   entity.Goods{},
		entities: entity.AllGoods{},
	}
}

// Create repository use database transaction
func (c *goodsRepository) Create(ctx context.Context, gds entity.Goods) (entity.Goods, error) {
	db := database.Instance
	// start transaction
	return gds, db.Transaction(func(tx *gorm.DB) error {

		gds.ID = utils.Uid(13)
		gds.Name = strings.ToTitle(gds.Name)

		cat := catEntity.Category{}
		cf := db.Where("id = ?", gds.CategoryId).First(&cat)

		logger.New(cf.Error, logger.SetType(logger.ERROR))

		// check if category not exist and let's create
		if cat.ID == "" {
			catId := utils.Uid(13)
			catR := catEntity.Category{Name: strings.ToTitle(gds.CategoryId), ID: catId}
			cr := db.Create(&catR)
			if cr.Error != nil {
				logger.New(cf.Error, logger.SetType(logger.ERROR))
				return cr.Error
			}
			gds.CategoryId = catId
			gds.Category = catR
		}

		r := db.WithContext(ctx).Create(&gds)

		// rollback if error
		if r.Error != nil {
			redis_client.Del(ctx, gds.ID)
			return r.Error
		}

		redis_client.Set(ctx, gds.ID, gds)
		redis_client.Del(ctx, "goods")

		return nil
	})
}

func (c *goodsRepository) Update(ctx context.Context, id string, gds entity.Goods) error {
	gds.Name = strings.ToTitle(gds.Name)

	d := gds
	d.ID = id
	redis_client.Set(ctx, id, d)
	redis_client.Del(ctx, "goods")
	return database.Instance.WithContext(ctx).Where("id = ?", id).Updates(&gds).Error
}

func (c *goodsRepository) Delete(ctx context.Context, id string) error {
	redis_client.Del(ctx, id)
	redis_client.Del(ctx, "goods")
	return database.Instance.WithContext(ctx).Where("id = ?", id).Delete(&c.entity).Error
}

func (c *goodsRepository) FindById(ctx context.Context, id string) (entity.Goods, error) {
	e := c.entity
	r := database.Instance.WithContext(ctx).Preload("Category").Where("id = ?", id).First(&e)
	if r.Error == nil {
		redis_client.FindSet(ctx, id, e)
	}
	return e, r.Error
}

func (c *goodsRepository) Stock(ctx context.Context, id string, n int) error {
	e := c.entity

	d := database.Instance.WithContext(ctx).Where("id = ?", id).First(&e)
	logger.New(d.Error, logger.SetType(logger.ERROR))

	redis_client.Set(ctx, id, e)
	redis_client.Del(ctx, "goods")
	r := database.Instance.WithContext(ctx).Preload("Category").Model(&entity.Goods{}).Where("id = ?", id).Update("stock", e.Stock+n)
	return r.Error
}

func (c *goodsRepository) All(ctx context.Context, page int, take int) (utils.Pagination, error) {
	e := c.entities
	var pagination utils.Pagination
	var totalRow int64

	r := database.Instance.WithContext(ctx).Model(&e)
	r.Count(&totalRow)
	r.Scopes(pagination.Paginate(page, take))
	r.Preload("Category").Find(&e)

	pagination.Rows = e
	pagination.CalculatePage(float64(totalRow))

	if totalRow > 0 {
		redis_client.FindSet(ctx, "goods", pagination)
	}
	return pagination, r.Error
}
