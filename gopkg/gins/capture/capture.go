package gins

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Capture 默认已支持mongo、json的序列化
type Capture struct {
	*filterResult
	Latency      time.Duration `json:"latency" bson:"latency"`
	Request      *Request      `json:"request" bson:"request"`
	Response     *Response     `json:"response" bson:"response"`
	CustomFields CustomFields  `json:"custom_fields" bson:"custom_fields"`
}

func NewCapture(ctx *gin.Context, request *Request, options Options) *Capture {
	startTime := time.Now()
	filterRes := newFilterResult(options, request)
	capture := &Capture{
		Request:      request,
		Response:     NewResponse(ctx, newWriterCaptureGinNext(ctx)),
		Latency:      time.Since(startTime),
		filterResult: filterRes,
		CustomFields: options.GetCustomField(ctx),
	}
	if filterRes != nil {
		return capture
	}

	response := capture.Response
	if response == nil {
		return capture
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return capture
	}

	return nil
}

func (c *Capture) GetCustomField(key string) string {
	return c.CustomFields.Get(key)
}

func newWriterCaptureGinNext(ctx *gin.Context) *CustomResponseWriter {
	writer := newCustomResponseWriter(ctx.Writer)
	ctx.Writer = writer
	ctx.Next()
	return writer
}
