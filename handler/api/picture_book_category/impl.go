package picture_book_category

import (
	"go-agent/gopkg/gins"
	"go-agent/internal/service"
	"go-agent/internal/service/picture_book_category"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	g                          *gin.RouterGroup
	pictureBookCategoryService service.PictureBookCategory
}

func NewHandler(g *gin.RouterGroup) gins.Handler {
	return &Handler{
		g:                          g,
		pictureBookCategoryService: picture_book_category.NewService(),
	}
}

func (h *Handler) RegisterRoutes() {
	g := h.g.Group("/picture/book/category")
	g.GET("/list", h.PagingPictureBookCategory)
}
