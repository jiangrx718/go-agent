package picture_book

import (
	"go-agent/gopkg/gorms"
	"go-agent/internal/dao"
	"go-agent/internal/model"
	"context"
)

func (d *Dao) Pagination(ctx context.Context, page gorms.Page) (*gorms.Paging[*model.SPictureBook], error) {
	paging, err := gorms.PaginationQuery(
		dao.SPictureBook.Order(
			dao.SPictureBook.Position.Desc(),
		).FindByPage, gorms.Page{
			PageIndex: page.PageIndex,
			PageSize:  page.PageSize,
		})
	if err != nil {
		return nil, d.ConvertError(err)
	}

	return paging, nil
}
