package chinese

import (
	"go-agent/gopkg/gins"
	"go-agent/handler/api/chinese/request"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Detail(ctx *gin.Context) {
	var req request.DetailRequest

	if err := ctx.Bind(&req); err != nil {
		gins.BadRequest(ctx, err)
		return
	}

	res, err := h.chineseService.Detail(ctx, req.Chinese)
	if err != nil {
		gins.ServerError(ctx, err)
		return
	}

	gins.StatusOK(ctx, res)
}
