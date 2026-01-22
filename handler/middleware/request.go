package middleware

import (
	capture "go-agent/gopkg/gins/capture"
	"go-agent/gopkg/utils"

	"github.com/gin-gonic/gin"
)

func RequestCapture() gin.HandlerFunc {
	return capture.RequestCapture(capture.Options{
		FilterPaths: []string{},
	}, func(ctx *gin.Context, request *capture.Request) {
		ctx.Set(utils.ClientIPKey, request.ClientIP)
	}, func(ctx *gin.Context, capture *capture.Capture) {
	})
}
