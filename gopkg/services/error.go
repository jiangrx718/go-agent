package services

import (
	"errors"
	"fmt"
)

var errorCodes = []Error{}

func GetErrorCodes() []Error {
	return errorCodes
}

type Error interface {
	Error() string
	GetCode() int
}

func asError(err error) Error {
	if err == nil {
		return nil
	}

	var codeErr Error
	if ok := errors.As(err, &codeErr); !ok {
		return nil
	}

	return codeErr
}

type BaseError struct {
	Code int    `json:"code"` // 业务异常码
	Msg  string `json:"msg"`  // 业务异常描述
}

func NewError(code int, msg string) *BaseError {
	baseError := &BaseError{
		Code: code,
		Msg:  msg,
	}
	errorCodes = append(errorCodes, baseError)
	return baseError
}

func (b *BaseError) GetCode() int {
	return b.Code
}

func (b *BaseError) Error() string {
	return b.Msg
}

func (b *BaseError) Sprintf(values ...any) *BaseError {
	return &BaseError{
		Code: b.Code,
		Msg:  fmt.Sprintf(b.Msg, values...),
	}
}
