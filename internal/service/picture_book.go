package service

import (
	"go-agent/gopkg/gorms"
	"go-agent/gopkg/services"
	"context"
)

type PictureBook interface {
	PagingPictureBook(ctx context.Context, page gorms.Page) (services.Result, error)
}
