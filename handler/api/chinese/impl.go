package chinese

import (
	"go-agent/gopkg/gins"
	"go-agent/internal/service"
	"go-agent/internal/service/chinese"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	g              *gin.RouterGroup
	chineseService service.Chinese
}

func NewHandler(g *gin.RouterGroup) gins.Handler {
	return &Handler{
		g:              g,
		chineseService: chinese.NewService(),
	}
}

func (h *Handler) RegisterRoutes() {
	g := h.g.Group("/chinese")
	g.POST("/detail", h.Detail)
}
