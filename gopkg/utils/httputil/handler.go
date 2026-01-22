package httputil

import (
	"net/http"

	"go-agent/gopkg/utils"
	"github.com/gin-gonic/gin"
)

type HttpError struct {
	Code int    `json:"code" example:"400"`
	Msg  string `json:"msg"  example:"参数错误"`
}

func BadRequest(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusOK, HttpError{
		Code: 400,
		Msg:  err.Error(),
	})
}

func ServerError(ctx *gin.Context, err error) {
	msg := err.Error()
	if utils.IsProduction() {
		msg = "服务器错误"
	}
	ctx.AbortWithStatusJSON(500, HttpError{
		Code: 500,
		Msg:  msg,
	})
}

func UnAuthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(200, HttpError{
		Code: 401,
		Msg:  "需要登录",
	})
}

func Forbidden(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(403, HttpError{
		Code: 403,
		Msg:  "需要登录",
	})
}

type Pagination struct {
	Limit  int64 `json:"limit" form:"limit,default=10" validate:"required"`
	Offset int64 `json:"offset" form:"offset"`
}
