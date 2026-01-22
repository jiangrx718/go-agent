package picture_book

import (
	"go-agent/gopkg/gorms"
	"go-agent/gopkg/log"
	"go-agent/gopkg/paging"
	"go-agent/gopkg/services"
	"go-agent/internal/model"
	"context"

	"go.uber.org/zap"
)

func (s *Service) PagingPictureBook(ctx context.Context, page gorms.Page) (services.Result, error) {
	logPrefix := "/internal/service/picture_book: Service.PagingPictureBook()"

	demoPaging, err := s.pictureBookDao.Pagination(ctx, page)
	if err != nil {
		log.Sugar().Error(logPrefix, zap.Any("picture_book dao pagination error", err), zap.Any("page", page))
		return services.Failed(ctx, err)
	}
	return services.Success(ctx, paging.NewPaging(demoPaging.Total, NewPictureBookS(demoPaging.List)))
}

func NewPictureBookS(demoEntities []*model.SPictureBook) []*model.SPictureBook {

	return demoEntities
}
