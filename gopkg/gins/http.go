package gins

import (
	"go-agent/gopkg/services"
	"go-agent/gopkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StatusOK(ctx *gin.Context, result services.Result) {
	ctx.AbortWithStatusJSON(http.StatusOK, result)
}

func BadRequest(ctx *gin.Context, err error) {
	StatusFailed(ctx, http.StatusBadRequest, err)
}

func ServerError(ctx *gin.Context, err error) {
	StatusFailed(ctx, http.StatusInternalServerError, err)
}

func Unauthorized(ctx *gin.Context) {
	StatusFailed(ctx, http.StatusUnauthorized, nil)
}

func StatusFailed(ctx *gin.Context, code int, err error) {
	if utils.IsProduction() || err == nil {
		ctx.AbortWithStatusJSON(code, services.NewResult(ctx, code, http.StatusText(code), nil))
		return
	}

	ctx.AbortWithStatusJSON(code, services.NewResult(ctx, code, err.Error(), nil))
}
