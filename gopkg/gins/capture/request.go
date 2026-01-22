package gins

import (
	"bytes"
	"io"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Method   string     `json:"method" bson:"method"`
	Path     string     `json:"path" bson:"path"`
	Query    url.Values `json:"query" bson:"query"`
	ClientIP string     `json:"client_ip" bson:"client_ip"`
	Body     any        `json:"body" bson:"body"`
}

func NewRequest(ctx *gin.Context) (*Request, error) {
	//log.Sugar().Info(ctx, "new request", zap.Any("x-forwarded-for", ctx.Request.Header.Get("x-forwarded-for")))
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}

	//log.Sugar().Info(ctx, "client_ip", zap.Any("ctx_client_ip", ctx.ClientIP()), zap.Any("ctx_remote_ip", ctx.RemoteIP()), zap.Any("x_forwarded_for", ctx.Request.Header.Get("x-forwarded-for")))
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return &Request{
		Method:   ctx.Request.Method,
		Path:     ctx.Request.URL.Path,
		Query:    ctx.Request.URL.Query(),
		ClientIP: ctx.ClientIP(),
		Body:     BytesToAny(body),
	}, nil
}
