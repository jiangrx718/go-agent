Agent 包说明

此包提供一个简单的 `Agent` 接口以及两种实现示例：

- 本地实现（默认）：`NewLocalAgent()`，仅返回输入的回声，便于本地开发与测试。
- langchaingo 实现（可选）：在启用 build tag `langchaingo` 时可用。该实现位于 `langchaingo_agent.go`，构建示例：

使用 langchaingo 示例：

环境变量：

- `OPENAI_API_KEY`：设置你的 OpenAI API Key
- 可选运行时开关 `USE_REMOTE_AGENT=1`（示例），但该开关需要配合 build tag 使用

构建并运行（示例）：

```
OPENAI_API_KEY=sk-... USE_REMOTE_AGENT=1 go run -tags langchaingo ./...
```

注意：`langchaingo` 的 API 可能会随版本变化，上述示例实现可能需要根据实际版本调整。
