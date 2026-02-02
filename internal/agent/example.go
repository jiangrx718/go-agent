package agent

import (
	"context"
	"fmt"
	"go-agent/gopkg/log"
	"os"
)

// NewAgentFromEnv 根据环境选择 Agent：默认返回本地实现。
// 若希望使用 langchaingo 实现，请在构建时启用 build tag `langchaingo`
// 并设置环境变量 OPENAI_API_KEY。示例：
//   OPENAI_API_KEY=... go run -tags langchaingo ./...
func NewAgentFromEnv() (Agent, error) {
    // 优先检查 Eino Agent
    if os.Getenv("USE_EINO_AGENT") == "1" {
        return NewEinoAgentFromEnv()
    }

    // 运行时开关：如果设置 USE_REMOTE_AGENT=1 且构建时启用了 langchaingo tag，
    // 可以返回真实实现。这里以简单方式说明流程。
    if os.Getenv("USE_REMOTE_AGENT") == "1" {
        // 当启用 langchaingo build tag 时，用户可以自己直接调用 NewLangChainAgent()
        // 这里不直接引用以避免在未启用 tag 时的编译问题。
        return NewLocalAgent(), fmt.Errorf("remote agent requested but langchaingo build tag not enabled; build with -tags langchaingo to enable")
    }
    return NewLocalAgent(), nil
}

// ExampleRun 展示如何从项目中调用 Agent
func ExampleRun(ctx context.Context, prompt string) {
    ag, err := NewAgentFromEnv()
    if err != nil {
        log.Sugar().Warnf("agent: %v, falling back to local agent", err)
    }

    resp, err := ag.Handle(ctx, prompt)
    if err != nil {
        log.Sugar().Errorf("agent handle error: %v", err)
        return
    }

    log.Sugar().Infof("agent response: %s", resp)
}
