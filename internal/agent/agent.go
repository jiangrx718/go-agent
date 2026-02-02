package agent

import (
	"context"

	"github.com/cloudwego/eino/schema"
)

// Agent 是一个简单的对话/任务代理接口，返回字符串结果或错误
type Agent interface {
	Handle(ctx context.Context, prompt string) (string, error)
}

// StreamAgent 扩展 Agent 接口以支持流式传输
type StreamAgent interface {
	Agent
	StreamHandle(ctx context.Context, prompt string) (*schema.StreamReader[*schema.Message], error)
}

// NewLocalAgent 返回一个不会调用外部 LLM 的本地实现（用于开发与测试）
func NewLocalAgent() Agent {
	return &localAgent{}
}

// NewEinoAgentFromEnv 返回一个使用 Eino 的实现
func NewEinoAgentFromEnv() (Agent, error) {
	return NewEinoAgent()
}
