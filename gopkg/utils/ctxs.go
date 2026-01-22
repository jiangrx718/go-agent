package utils

import (
	"context"
)

const (
	RequestIDKey = "x-request-id"
	ClientIPKey  = "client-ip"
)

func GetRequestID(ctx context.Context) string {
	return GetString(ctx, RequestIDKey)
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func SetClientIP(ctx context.Context, clientIP string) context.Context {
	return context.WithValue(ctx, ClientIPKey, clientIP)
}

func GetClientIP(ctx context.Context) string {
	return GetString(ctx, ClientIPKey)
}

func GetString(ctx context.Context, key string) string {
	valueAny := ctx.Value(key)
	if valueAny == nil {
		return ""
	}

	valueString, ok := valueAny.(string)
	if !ok {
		return ""
	}

	return valueString
}
