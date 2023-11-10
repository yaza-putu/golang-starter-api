package service

import (
	"context"
	"github.com/yaza-putu/golang-starter-api/src/app/category/entity"
	"github.com/yaza-putu/golang-starter-api/src/app/category/repository"
	"github.com/yaza-putu/golang-starter-api/src/app/category/validation"
	"github.com/yaza-putu/golang-starter-api/src/http/response"
	"github.com/yaza-putu/golang-starter-api/src/logger"
)

type categoryService struct {
	repository repository.CategoryInterface
}

type optFunc func(*repository.CategoryInterface)

// NewCategoryService with optional parsing repository
// benefit to use without parsing repository on call and can parsing repository for make mock test
func NewCategoryService(opts ...optFunc) CategoryInterface {
	o := repository.NewCategory()

	for _, fn := range opts {
		fn(&o)
	}

	return &categoryService{
		repository: o,
	}
}

func Mock(s repository.CategoryInterface) optFunc {
	return func(categoryInterface *repository.CategoryInterface) {
		categoryInterface = &s
	}
}

func (c *categoryService) All(ctx context.Context, page int, take int) response.DataApi {
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

func (c *categoryService) Create(ctx context.Context, cat validation.Category) response.DataApi {
	rc := make(chan response.DataApi)

	go func() {
		defer close(rc)
		r, err := c.repository.Create(ctx, entity.Category{
			Name: cat.Name,
		})

		if err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
			rc <- response.Api(response.SetCode(500), response.SetMessage(err))
		}

		rc <- response.Api(response.SetCode(200), response.SetData(r), response.SetMessage("Create category successfully"))
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

func (c *categoryService) Update(ctx context.Context, id string, cat validation.CategoryWithIgnore) response.DataApi {
	rc := make(chan response.DataApi)

	go func() {
		err := c.repository.Update(ctx, id, entity.Category{Name: cat.Name})
		defer close(rc)

		if err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
			rc <- response.Api(response.SetCode(500), response.SetMessage(err))
		}

		rc <- response.Api(response.SetCode(200), response.SetMessage("Update category successfully"))
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

func (c *categoryService) Delete(ctx context.Context, id string) response.DataApi {
	rc := make(chan response.DataApi)
	go func() {
		err := c.repository.Delete(ctx, id)
		defer close(rc)

		if err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
			rc <- response.Api(response.SetCode(500), response.SetMessage(err))
		}

		rc <- response.Api(response.SetCode(200), response.SetMessage("Delete category successfuly"))
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

func (c *categoryService) FindById(ctx context.Context, id string) response.DataApi {
	rc := make(chan response.DataApi)

	go func() {
		r, err := c.repository.FindById(ctx, id)
		defer close(rc)

		if err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
			rc <- response.Api(response.SetCode(404), response.SetStatus(false), response.SetMessage(err))
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
