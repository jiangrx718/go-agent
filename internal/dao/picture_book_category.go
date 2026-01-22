package dao

import (
	"go-agent/gopkg/gorms"
	"go-agent/internal/model"
	"context"
)

type PictureBookCategory interface {
	Pagination(ctx context.Context, page gorms.Page) (*gorms.Paging[*model.SPictureBookCategory], error)
}
