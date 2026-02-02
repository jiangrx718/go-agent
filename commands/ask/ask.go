package ask

import (
	"context"
	"encoding/json"
	"fmt"
	"go-agent/internal/agent"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

type ollamaTagsResponse struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "ask",
		Usage: "向 LLM 代理提问",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "指定要使用的 Ollama 模型",
			},
			&cli.StringFlag{
				Name:  "base-url",
				Usage: "指定 Ollama 服务地址",
				Value: "http://localhost:11434",
			},
		},
		Action: func(c *cli.Context) error {
			prompt := c.Args().First()
			if prompt == "" {
				return fmt.Errorf("请提供提示词")
			}

			// 设置 Base URL
			baseURL := c.String("base-url")
			// eino agent 内部会自动追加 /v1，但这里我们先不加，因为 detectFirstModel 需要原始地址
			// 传递给 Eino Agent 的环境变量 OLLAMA_BASE_URL 应该包含 /v1 吗？
			// Eino agent 内部实现默认是 http://localhost:11434/v1
			// 为了保持一致性，我们在 ask.go 中只处理 host:port，在传递给 agent 时再拼接 /v1

			// 如果 baseURL 没有 /v1 后缀，且 agent 需要 /v1，我们得小心。
			// Eino agent 代码里:
			// if baseURL == "" { baseURL = "http://localhost:11434/v1" }
			//
			// 如果我们在 ask.go 设置了 OLLAMA_BASE_URL="http://localhost:11434"，agent 会直接使用它，
			// 导致请求发往 http://localhost:11434/chat/completions 而不是 http://localhost:11434/v1/chat/completions

			if baseURL != "" {
				// 确保 OLLAMA_BASE_URL 包含 /v1
				// 但 detectFirstModel 需要原始地址

				// 1. 设置环境变量供 Agent 使用 (加上 /v1)
				agentBaseURL := baseURL
				// 简单判断是否已有 v1
				// 更好的做法可能是让 Agent 内部处理，但为了快速修复，我们在外部处理
				os.Setenv("OLLAMA_BASE_URL", agentBaseURL+"/v1")
			} else {
				baseURL = "http://localhost:11434"
				// 默认情况 agent 内部会处理 /v1
			}

			// 确定模型
			modelName := c.String("model")
			if modelName == "" {
				modelName = os.Getenv("OLLAMA_MODEL")
			}

			// 如果仍未指定模型，尝试自动检测
			if modelName == "" {
				fmt.Println("未指定模型，正在检测本地 Ollama 模型...")
				detectedModel, err := detectFirstModel(baseURL)
				if err == nil && detectedModel != "" {
					modelName = detectedModel
					fmt.Printf("自动检测到模型: %s\n", modelName)
				} else {
					fmt.Printf("自动检测模型失败: %v\n", err)
					// 只有在检测失败时才使用硬编码默认值，或者让 Agent 内部处理（Agent 内部默认为 llama3）
				}
			}

			if modelName != "" {
				os.Setenv("OLLAMA_MODEL", modelName)
			}

			// 确保使用 Eino Agent
			os.Setenv("USE_EINO_AGENT", "1")

			ag, err := agent.NewAgentFromEnv()
			if err != nil {
				return err
			}

			fmt.Printf("正在向 Agent 提问 (使用 Eino + Ollama, 模型: %s): %s\n", os.Getenv("OLLAMA_MODEL"), prompt)
			resp, err := ag.Handle(context.Background(), prompt)
			if err != nil {
				return err
			}

			fmt.Println("\n回答:")
			fmt.Println(resp)
			return nil
		},
	}
}

func detectFirstModel(baseURL string) (string, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(baseURL + "/api/tags")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama api returned status: %d", resp.StatusCode)
	}

	var tagsResp ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return "", err
	}

	if len(tagsResp.Models) > 0 {
		return tagsResp.Models[0].Name, nil
	}

	return "", fmt.Errorf("no models found in ollama")
}
