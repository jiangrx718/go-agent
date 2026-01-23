package agent

import "context"

// Agent 是一个简单的对话/任务代理接口，返回字符串结果或错误
type Agent interface {
    Handle(ctx context.Context, prompt string) (string, error)
}

// NewLocalAgent 返回一个不会调用外部 LLM 的本地实现（用于开发与测试）
func NewLocalAgent() Agent {
    return &localAgent{}
}
