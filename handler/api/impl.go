package api

import (
	"go-agent/gopkg/gins"
	"go-agent/handler/api/chinese"
	"go-agent/handler/api/picture_book"
	"go-agent/handler/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine *gin.Engine
}

func NewHandler(engine *gin.Engine) gins.Handler {
	return &Handler{
		engine: engine,
	}
}

func (h *Handler) RegisterRoutes() {
	config := cors.DefaultConfig()
	config.AllowHeaders = append([]string{}, config.AllowHeaders...)
	config.AllowAllOrigins = true
	h.engine.Use(cors.New(config))

	g := h.engine.Group("/api", middleware.RequestCapture())
	handlers := []gins.Handler{
		picture_book.NewHandler(g),
		chinese.NewHandler(g),
	}

	for _, handler := range handlers {
		handler.RegisterRoutes()
	}
}
