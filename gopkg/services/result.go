package services

import (
	"go-agent/gopkg/utils"
	"context"
	"reflect"
)

type Result interface {
	GetCode() int
	GetMsg() string
	GetData() any
}

func Failed(ctx context.Context, err error) (Result, error) {
	codeError := asError(err)
	if codeError == nil {
		return nil, err
	}

	return newBaseResult(ctx, codeError.GetCode(), codeError.Error(), nil), nil
}

func Success(ctx context.Context, data any) (Result, error) {
	return newBaseResult(ctx, successfulDefaultCode, successfulDefaultMessage, data), nil
}

func NewResult(ctx context.Context, code int, msg string, data any) Result {
	return newBaseResult(ctx, code, msg, data)
}

type BaseResult struct {
	Code      int    `json:"code"`       // 响应码: 0=成功, 1=失败
	Msg       string `json:"msg"`        // 响应消息
	Data      any    `json:"data"`       // 响应体
	RequestID string `json:"request_id"` // 请求ID
}

func newBaseResult(ctx context.Context, code int, message string, data any) *BaseResult {
	v := reflect.ValueOf(data)
	if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		data = make(map[string]any)
	}

	return &BaseResult{
		Code:      code,
		Msg:       message,
		Data:      data,
		RequestID: utils.GetRequestID(ctx),
	}
}

func (b *BaseResult) GetCode() int {
	return b.Code
}

func (b *BaseResult) GetMsg() string {
	return b.Msg
}

func (b *BaseResult) GetData() any {
	return b.Data
}
