package db

import (
	"gorm.io/gorm"
	"math"
)

type PaginationValidation struct {
	Page int `query:"page" validate:"required"`
	Take int `query:"take" validate:"required"`
}

type Pagination struct {
	Limit      int    `json:"limit,omitempty;query:limit"`
	Page       int    `json:"page,omitempty;query:page"`
	Sort       string `json:"sort,omitempty;query:sort"`
	TotalRows  uint64 `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	Rows       any    `json:"rows"`
}

func (p *Pagination) SetOffset(offset int) {
	p.Limit = offset
}

func (p *Pagination) SetSort(sort string) {
	p.Sort = sort
}

func (p *Pagination) SetPage(page int) {
	p.Page = page
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}

func (p *Pagination) Paginate(pg int, of int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := pg

		if page <= 0 {
			page = 1
		}

		p.SetPage(page)

		pageSize := of
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		p.SetOffset(pageSize)

		return db.Offset(p.GetOffset()).Limit(p.GetLimit()).Order(p.GetSort())
	}
}

func (p *Pagination) CalculatePage(totalRow float64) {
	totalPage := math.Round(totalRow / float64(p.Limit))

	if totalPage < 1 && totalRow > 0 {
		totalPage = 1
	}

	p.TotalRows = uint64(totalRow)
	p.TotalPages = int(totalPage)

}
