package request

type DetailRequest struct {
	Chinese string `json:"chinese" binding:"required"`
}
