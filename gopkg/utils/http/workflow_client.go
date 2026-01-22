package httputil

import (
	"go-agent/gopkg/log"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// TaskType 定义工作流任务类型
type TaskType string

const (
	TaskTranslate TaskType = "translate"
	TaskSummarize TaskType = "summarize"
	TaskExplain   TaskType = "explain"
)

// LLMConfig 语言模型配置
type LLMConfig struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
}

// EncryptionConfig 加密配置
type EncryptionConfig struct {
	Type   string `json:"type,omitempty"`
	IsOpen bool   `json:"is_open"`
}

// WorkflowParams 通用入口参数
// 根据不同 task_type 使用 Query/TargetLanguage/FullText 等字段
type WorkflowParams struct {
	Query          string           // 必填：translate/summarize/explain 的查询文本
	TargetLanguage string           // translate 目标语言
	FullText       string           // summarize/explain 全文或上下文
	Temperature    float64          // 采样温度，可选
	MaxTokens      int              // 最大生成 token，可选
	LLM            LLMConfig        // LLM 配置
	Encryption     EncryptionConfig // 加密配置
}

// GeneralRequest 请求体结构（与接口契合）
type GeneralRequest struct {
	WorkflowID string `json:"workflow_id"`
	Inputs     struct {
		TextProcessing struct {
			Query          string    `json:"query"`
			TaskType       TaskType  `json:"task_type"`
			TargetLanguage string    `json:"target_language,omitempty"`
			FullText       string    `json:"full_text,omitempty"`
			LLMConfig      LLMConfig `json:"llm_config"`
			Temperature    float64   `json:"temperature,omitempty"`
			MaxTokens      int       `json:"max_tokens,omitempty"`
		} `json:"text_processing"`
	} `json:"inputs"`
	Configs struct {
		EncryptionConfig EncryptionConfig `json:"encryption_config"`
	} `json:"configs"`
}

// GeneralResponse 响应体结构（与示例兼容）
type GeneralResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		WorkflowID    *string        `json:"workflow_id"`
		WorkflowName  *string        `json:"workflow_name"`
		TaskID        *string        `json:"task_id"`
		CreateTime    float64        `json:"create_time"`
		CompletedTime float64        `json:"completed_time"`
		Results       map[string]any `json:"results"`
	} `json:"data"`
}

// 默认高性能 http.Client（连接复用、合理超时）
var defaultClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
	Timeout: 30 * time.Second,
}

// RunGeneralWorkflow 发送通用工作流请求
// - endpoint 例如: "http://127.0.0.1:8080/v1/openapi/workflow/running/general"
// - workflowID 固定 UUID（示例中为 96158309-d167-545a-af18-59f995d8ec8d）
// - params 根据 taskType 填充 Query/TargetLanguage/FullText 等
// - client 可传入自定义 *http.Client；为 nil 时使用默认 client
func RunGeneralWorkflow(ctx context.Context, endpoint string, req GeneralRequest, client *http.Client) (*GeneralResponse, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint is empty")
	}
	if req.WorkflowID == "" {
		return nil, errors.New("workflow_id is empty")
	}
	if client == nil {
		client = defaultClient
	}

	// 编码 JSON（使用 bytes.Buffer 避免多次分配）
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, err
	}

	// 打印请求 URL 与参数（敏感字段脱敏）
	func() {
		// 复制一份用于日志，避免修改真实请求体
		logReq := req
		// 使用结构化字段打印，避免转义引号
		log.Logger().Info(
			"RunGeneralWorkflow request",
			zap.String("url", endpoint),
			zap.Any("body", logReq),
		)
	}()

	// 构建请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解码响应（无论状态码，尝试按约定格式解析）
	var gr GeneralResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return nil, err
	}
	// 当服务返回非 2xx 时，仍将解析结果返回给调用方以便处理 message/结果
	return &gr, nil
}

// 便捷包装：固定默认 endpoint 与 workflowID
const DefaultEndpoint = "http://10.122.9.26:31373/v1/openapi/workflow/running/general"
const DefaultWorkflowID = "96158309-d167-545a-af18-59f995d8ec8d"

// RunGeneralWorkflowDefault 使用默认 endpoint 与 workflowID
func RunGeneralWorkflowDefault(ctx context.Context, req GeneralRequest) (*GeneralResponse, error) {
	return RunGeneralWorkflow(ctx, DefaultEndpoint, req, nil)
}

// 使用示例：
// 1) 翻译（默认 endpoint/workflowID）：
// resp, err := httputil.RunGeneralWorkflowDefault(ctx, httputil.TaskTranslate, httputil.WorkflowParams{
//   Query:          "Artificial Intelligence is transforming the world.",
//   TargetLanguage: "中文",
//   Temperature:    0.3,
//   MaxTokens:      1024,
//   LLM: httputil.LLMConfig{
//     Provider: "openai_api_compatible",
//     Model:    "Qwen2.5-VL-7B-Instruct",
//     BaseURL:  "https://studio.bd.kxsz.net:9443/v1",
//     APIKey:   "sk-xxx",
//   },
//   Encryption: httputil.EncryptionConfig{Type: "PKCS1_OAEP", IsOpen: false},
// })
// // 日志会自动打印请求 URL 和参数（APIKey 已脱敏）
// if err != nil { /* handle */ }
// if resp != nil && resp.Code == 200 {
//   _ = resp.Data.Results["result"]
// }

// 2) 总结全文（自定义 endpoint/workflowID）：
// resp, err := httputil.RunGeneralWorkflow(
//   ctx,
//   "http://127.0.0.1:8080/v1/openapi/workflow/running/general",
//   "96158309-d167-545a-af18-59f995d8ec8d",
//   httputil.TaskSummarize,
//   httputil.WorkflowParams{
//     FullText:    "<长文本>...",
//     Temperature: 0.2,
//     MaxTokens:   1536,
//     LLM: httputil.LLMConfig{
//       Provider: "openai_api_compatible",
//       Model:    "Qwen2.5-32B-Instruct",
//       BaseURL:  "https://studio.bd.kxsz.net:9443/v1",
//       APIKey:   "sk-xxx",
//     },
//     Encryption: httputil.EncryptionConfig{Type: "PKCS1_OAEP", IsOpen: false},
//   },
//   nil,
// )
// if err != nil { /* handle */ }
// if resp != nil && resp.Code == 200 {
//   _ = resp.Data.Results["result"]
// }

// 3) 解释/讲解（传入 Query 提示词）：
// resp, err := httputil.RunGeneralWorkflowDefault(ctx, httputil.TaskExplain, httputil.WorkflowParams{
//   Query:       "Explain what is a vector database and its use cases.",
//   Temperature: 0.7,
//   LLM: httputil.LLMConfig{
//     Provider: "openai_api_compatible",
//     Model:    "Qwen2.5-7B-Instruct",
//     BaseURL:  "https://studio.bd.kxsz.net:9443/v1",
//     APIKey:   "sk-xxx",
//   },
// })
// if err != nil { /* handle */ }
// if resp != nil && resp.Code == 200 {
//   _ = resp.Data.Results["result"]
// }

// 4) 段落总结（将段落作为 FullText 或 Query 均可，推荐 FullText）：
// paragraph := "This is a single paragraph that needs summarization. ..."
// resp, err := httputil.RunGeneralWorkflowDefault(ctx, httputil.TaskSummarize, httputil.WorkflowParams{
//   FullText:    paragraph,
//   Temperature: 0.3,
//   LLM: httputil.LLMConfig{
//     Provider: "openai_api_compatible",
//     Model:    "Qwen2.5-7B-Instruct",
//     BaseURL:  "https://studio.bd.kxsz.net:9443/v1",
//     APIKey:   "sk-xxx",
//   },
// })
// if err != nil { /* handle */ }
// if resp != nil && resp.Code == 200 {
//   _ = resp.Data.Results["result"]
// }
