package service

import (
	"context"
	"github.com/yaza-putu/golang-starter-api/src/app/goods/entity"
	"github.com/yaza-putu/golang-starter-api/src/app/goods/repository"
	"github.com/yaza-putu/golang-starter-api/src/app/goods/validation"
	"github.com/yaza-putu/golang-starter-api/src/http/response"
	"github.com/yaza-putu/golang-starter-api/src/logger"
)

type goodsService struct {
	repository repository.GoodsInterface
}

type optFunc func(*repository.GoodsInterface)

// NewGoods with optional parsing repository
// benefit to use without parsing repository on call and can parsing repository for make mock test
func NewGoods(opts ...optFunc) GoodInterface {
	o := repository.NewGoods()

	for _, fn := range opts {
		fn(&o)
	}

	return &goodsService{
		repository: o,
	}
}

func Mock(s repository.GoodsInterface) optFunc {
	return func(categoryInterface *repository.GoodsInterface) {
		categoryInterface = &s
	}
}

func (c *goodsService) All(ctx context.Context, page int, take int) response.DataApi {
	rc := make(chan response.DataApi)

	go func() {
		r, err := c.repository.All(ctx, page, take)
		if err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
			rc <- response.Api(response.SetCode(500), response.SetMessage(err))
		}

		if r.TotalRows == 0 {
			rc <- response.Api(response.SetData([]string{}))
		}

		rc <- response.Api(response.SetCode(200), response.SetData(r))
		defer close(rc)
	}()

	for {
		select {
		case <-ctx.Done():
			return response.TimeOut()
		case res := <-rc:
			return res
		}
	}
}

func (c *goodsService) Create(ctx context.Context, gds validation.Goods) response.DataApi {
	rc := make(chan response.DataApi)

	go func() {
		defer close(rc)
		r, err := c.repository.Create(ctx, entity.Goods{
			Name:       gds.Name,
			CategoryId: gds.CategoryId,
		})

		if err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
			rc <- response.Api(response.SetCode(500), response.SetMessage(err))
		}

		rc <- response.Api(response.SetCode(200), response.SetData(r), response.SetMessage("Create goods successfully"))
	}()

	for {
		select {
		case <-ctx.Done():
			return response.TimeOut()
		case res := <-rc:
			return res
		}
	}
}

func (c *goodsService) Update(ctx context.Context, id string, gds validation.Goods) response.DataApi {
	rc := make(chan response.DataApi)

	go func() {
		err := c.repository.Update(ctx, id, entity.Goods{Name: gds.Name, CategoryId: gds.CategoryId})
		defer close(rc)

		if err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
			rc <- response.Api(response.SetCode(500), response.SetMessage(err))
		}

		rc <- response.Api(response.SetCode(200), response.SetMessage("Update goods successfully"))
	}()

	for {
		select {
		case <-ctx.Done():
			return response.TimeOut()
		case res := <-rc:
			return res
		}
	}
}

func (c *goodsService) Stock(ctx context.Context, id string, n int) response.DataApi {
	rc := make(chan response.DataApi)

	go func() {
		err := c.repository.Stock(ctx, id, n)
		defer close(rc)

		if err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
			rc <- response.Api(response.SetCode(500), response.SetMessage(err))
		}

		rc <- response.Api(response.SetCode(200), response.SetMessage("Update Stock successfully"))
	}()

	for {
		select {
		case <-ctx.Done():
			return response.TimeOut()
		case res := <-rc:
			return res
		}
	}
}

func (c *goodsService) Delete(ctx context.Context, id string) response.DataApi {
	rc := make(chan response.DataApi)
	go func() {
		err := c.repository.Delete(ctx, id)
		defer close(rc)

		if err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
			rc <- response.Api(response.SetCode(500), response.SetMessage(err))
		}

		rc <- response.Api(response.SetCode(200), response.SetMessage("Delete goods successfuly"))
	}()

	for {
		select {
		case <-ctx.Done():
			return response.TimeOut()
		case res := <-rc:
			return res
		}
	}
}

func (c *goodsService) FindById(ctx context.Context, id string) response.DataApi {
	rc := make(chan response.DataApi)

	go func() {
		r, err := c.repository.FindById(ctx, id)
		defer close(rc)

		if err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
			rc <- response.Api(response.SetCode(404), response.SetMessage(err))
		}

		rc <- response.Api(response.SetCode(200), response.SetData(r))
	}()

	for {
		select {
		case <-ctx.Done():
			return response.TimeOut()
		case res := <-rc:
			return res
		}
	}
}
