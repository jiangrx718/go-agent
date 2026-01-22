package service

import (
	"go-agent/gopkg/gorms"
	"go-agent/gopkg/services"
	"context"
)

type PictureBookCategory interface {
	PagingPictureBookCategory(ctx context.Context, page gorms.Page) (services.Result, error)
}
