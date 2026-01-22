package gins

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// CustomResponseWriter 自定义 ResponseWriter
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (c *CustomResponseWriter) Write(b []byte) (int, error) {
	c.body.Write(b)
	return c.ResponseWriter.Write(b)
}

func newCustomResponseWriter(writer gin.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{
		body:           bytes.NewBufferString(""),
		ResponseWriter: writer,
	}
}
