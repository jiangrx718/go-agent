package gins

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int    `json:"status_code" bson:"status_code"`
	Body       any    `json:"body" bson:"body"`
	Error      string `json:"error" bson:"error"`
}

func NewResponse(ctx *gin.Context, writer *CustomResponseWriter) *Response {
	return &Response{
		StatusCode: ctx.Writer.Status(),
		Body:       BytesToAny(writer.body.Bytes()),
		Error:      ctx.Errors.ByType(gin.ErrorTypePrivate).String(),
	}
}
