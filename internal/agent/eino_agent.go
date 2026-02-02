package agent

import (
	"context"
	"go-agent/gopkg/log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type EinoAgent struct {
	runnable compose.Runnable[[]*schema.Message, *schema.Message]
}

func NewEinoAgent() (*EinoAgent, error) {
	// 默认为本地 Ollama 实例
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	if baseURL == "" {
		// OpenAI SDK 通常会自动添加 /v1 后缀，或者根据客户端需求期望基础路径不含 /v1
		// 但在这里，我们使用的是 cloudwego/eino-ext/components/model/openai，它使用的是 meguminnnnnnnnn/go-openai
		// 对于 Ollama 来说，通常应该是 http://localhost:11434/v1 的形式
		baseURL = "http://localhost:11434/v1"
	}
	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "llama3"
	}

	// 创建 OpenAI 聊天模型（指向 Ollama）
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		BaseURL: baseURL,
		APIKey:  "ollama", // Ollama 不需要真实的 key，但某些客户端强制要求
		Model:   modelName,
	})
	if err != nil {
		return nil, err
	}

	// 创建一个简单的链：输入 -> 模型 -> 输出
	chain := compose.NewChain[[]*schema.Message, *schema.Message]()
	chain.AppendChatModel(chatModel)

	// 编译
	runnable, err := chain.Compile(context.Background())
	if err != nil {
		return nil, err
	}

	return &EinoAgent{
		runnable: runnable,
	}, nil
}

// Handle 实现 Agent 接口
func (a *EinoAgent) Handle(ctx context.Context, prompt string) (string, error) {
	// 将 prompt 转换为 Eino 消息
	msgs := []*schema.Message{
		schema.UserMessage(prompt),
	}

	// 生成
	resp, err := a.runnable.Invoke(ctx, msgs)

	if err != nil {
		log.Sugar().Errorf("eino agent invoke error: %v", err)
		return "", err
	}

	return resp.Content, nil
}

// StreamHandle 实现流式处理
func (a *EinoAgent) StreamHandle(ctx context.Context, prompt string) (*schema.StreamReader[*schema.Message], error) {
	// 将 prompt 转换为 Eino 消息
	msgs := []*schema.Message{
		schema.UserMessage(prompt),
	}

	// 生成流
	stream, err := a.runnable.Stream(ctx, msgs)
	if err != nil {
		log.Sugar().Errorf("eino agent stream error: %v", err)
		return nil, err
	}

	return stream, nil
}
