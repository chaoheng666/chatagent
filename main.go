// main.go —— 官方 ChatModel guide 示例 + DeepSeek 适配（2025年5月确认）
package main

import (
	"bufio"
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

// ChatRequest 前端发送的请求体
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse 返回给前端的响应
type ChatResponse struct {
	Reply string `json:"reply"`
}

func main() {
	// 尝试加载本地 api.env（如果存在）
	_ = loadDotEnv("api.env")

	ctx := context.Background()

	// 读取 API Key（支持两种环境变量名）
	apiKey := os.Getenv("OPENAPI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	cfg := &openai.ChatModelConfig{
		APIKey: apiKey,
		Model:  "gpt-5-nano",
	}
	if bu := os.Getenv("OPENAPI_BASE_URL"); bu != "" {
		cfg.BaseURL = bu
	}

	model, err := openai.NewChatModel(ctx, cfg)
	if err != nil {
		panic("模型创建失败: " + err.Error())
	}

	r := gin.Default()

	// 前端静态文件
	r.StaticFile("/", "static/index.html")
	r.Static("/static/", "static")

	// API 路由：接收前端消息并调用模型
	r.POST("/api/chat", func(c *gin.Context) {
		var req ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		messages := []*schema.Message{
			{Role: schema.System, Content: "你是一个友好、有帮助的中文助手，回答带点表情。"},
			{Role: schema.User, Content: req.Message},
		}

		resp, err := model.Generate(ctx, messages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		reply := ""
		if resp != nil {
			reply = resp.Content
		}
		c.JSON(http.StatusOK, ChatResponse{Reply: reply})
	})

	// 启动服务，监听在 8080
	r.Run(":8080")
}

// loadDotEnv 解析简单的 .env 文件样式，支持几种常见写法：
// - export KEY=VAL
// - export KEY = VAL
// - KEY=VAL
// 注：不会覆盖已存在的环境变量（仅在目标不存在时设置）。
func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, "'\"")
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
