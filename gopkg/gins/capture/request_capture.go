package gins

import (
	"go-agent/gopkg/gins"
	"github.com/gin-gonic/gin"
)

func RequestCapture(options Options, before func(ctx *gin.Context, request *Request), after func(ctx *gin.Context, capture *Capture)) gin.HandlerFunc {
	options.format()
	return func(ctx *gin.Context) {
		request, err := NewRequest(ctx)
		if err != nil {
			gins.ServerError(ctx, err)
			return
		}

		before(ctx, request)
		capture := NewCapture(ctx, request, options)
		if capture == nil {
			return
		}

		after(ctx, capture)
	}
}
