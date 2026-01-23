//go:build langchaingo
// +build langchaingo

package agent

import (
	"context"
	"errors"
	"go-agent/gopkg/log"
	"os"

	// 注意：此文件仅作为示例实现，受 build tag 控制（langchaingo）。
	// 使用前请确保在构建时启用该 tag：
	//    go build -tags langchaingo ./...
	// 同时需要设置环境变量 OPENAI_API_KEY

	"github.com/tmc/langchaingo"
	"github.com/tmc/langchaingo/llms/openai"
)

// LangChainAgent 使用 github.com/tmc/langchaingo 调用远端模型。
type LangChainAgent struct {
	client *langchaingo.Client
}

// NewLangChainAgent 从环境变量读取 API Key 并创建 Agent。
func NewLangChainAgent() (*LangChainAgent, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY is not set")
	}

	// 下面的构造基于 langchaingo 的常见用法；如有 API 变动请调整。
	llm := openai.New(openai.WithAPIKey(apiKey))
	client := langchaingo.New(llm)

	return &LangChainAgent{client: client}, nil
}

func (a *LangChainAgent) Handle(ctx context.Context, prompt string) (string, error) {
	if a == nil || a.client == nil {
		return "", errors.New("langchain client not initialized")
	}

	log.Sugar().Infof("agent: sending prompt to langchaingo: %s", prompt)

	// 示范调用，具体方法名请根据实际 langchaingo 版本调整
	resp, err := a.client.Call(ctx, prompt)
	if err != nil {
		log.Sugar().Errorf("agent: langchaingo call error: %v", err)
		return "", err
	}

	return resp, nil
}
