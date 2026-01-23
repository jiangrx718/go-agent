package agent

import (
	"context"
)

// localAgent 是一个非常简单的实现，仅用于本地测试与开发
type localAgent struct{}

func (a *localAgent) Handle(ctx context.Context, prompt string) (string, error) {
    // 简单回声实现 — 在没有外部模型或测试环境下可立即使用
    return "Echo: " + prompt, nil
}
