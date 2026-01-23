package service

import (
	"context"
	"go-agent/gopkg/services"
)

type Chinese interface {
	Detail(ctx context.Context, chinese string) (services.Result, error)
}
