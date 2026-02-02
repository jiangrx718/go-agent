package agent

import (
	"fmt"
	"go-agent/gopkg/gins"
	"go-agent/internal/agent"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	g *gin.RouterGroup
}

func NewHandler(g *gin.RouterGroup) gins.Handler {
	return &Handler{
		g: g,
	}
}

func (h *Handler) RegisterRoutes() {
	g := h.g.Group("/agent")
	g.GET("/chat", h.Chat)
}

// ChatRequest 请求结构
type ChatRequest struct {
	Prompt string `form:"prompt" binding:"required"`
	Model  string `form:"model"`
}

// Chat 流式问答接口
func (h *Handler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		gins.BadRequest(c, err)
		return
	}

	// 设置模型
	if req.Model != "" {
		os.Setenv("OLLAMA_MODEL", req.Model)
	}

	// 强制使用 Eino Agent
	os.Setenv("USE_EINO_AGENT", "1")
	// 确保 BaseURL 正确（根据之前的修复，Ollama 默认需要 /v1）
	if os.Getenv("OLLAMA_BASE_URL") == "" {
		os.Setenv("OLLAMA_BASE_URL", "http://localhost:11434/v1")
	}

	ag, err := agent.NewEinoAgent()
	if err != nil {
		gins.ServerError(c, fmt.Errorf("failed to create agent: %v", err))
		return
	}

	// 调用流式接口
	stream, err := ag.StreamHandle(c.Request.Context(), req.Prompt)
	if err != nil {
		gins.ServerError(c, fmt.Errorf("failed to start stream: %v", err))
		return
	}
	defer stream.Close()

	// 设置 SSE 头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		chunk, err := stream.Recv()
		if err == io.EOF {
			return false
		}
		if err != nil {
			// 如果流出错，尝试发送错误信息（虽然这时可能已经发送了部分数据）
			c.SSEvent("error", err.Error())
			return false
		}

		// 发送数据块
		c.SSEvent("message", chunk.Content)
		return true
	})
}
