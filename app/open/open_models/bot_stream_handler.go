package open_models

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// MessageFormat 消息格式类型（对标飞书/钉钉）
type MessageFormat string

const (
	FormatText     MessageFormat = "text"     // 纯文本
	FormatMarkdown MessageFormat = "markdown" // Markdown（最常用）
	FormatRichText MessageFormat = "richtext" // 富文本（JSON 结构）
	FormatHTML     MessageFormat = "html"     // HTML
)

// BotChunk Bot 流式消息片段
type BotChunk struct {
	Type     string                 `json:"type"`               // text/markdown/richtext/html
	Content  string                 `json:"content"`            // 内容
	Metadata map[string]interface{} `json:"metadata,omitempty"` // 元数据（@、附件等）
}

// BotStreamHandler 大厂标准 Bot 流式处理器
type BotStreamHandler struct {
	WebhookURL string
	Secret     string
	Timeout    int
	Payload    map[string]interface{}

	// 回调函数
	OnChunk    func(chunk *BotChunk) error // 每个片段
	OnComplete func() error                // 完成
	OnError    func(err error) error       // 错误
}

// ExecuteStreamCall 执行流式调用（支持多格式）
func (h *BotStreamHandler) ExecuteStreamCall(ctx context.Context) error {
	startTime := time.Now()

	// 1. 序列化 payload
	body, err := json.Marshal(h.Payload)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	// 2. 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", h.WebhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream, application/json")
	req.Header.Set("User-Agent", "Beaver-OpenPlatform/1.0")

	// 3. 添加签名（HMAC-SHA256）
	if h.Secret != "" {
		signature := generateHMACSignature(body, h.Secret, startTime.Unix())
		req.Header.Set("X-Webhook-Signature", signature)
		req.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", startTime.Unix()))
	}

	// 4. 发送请求
	timeout := h.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		if h.OnError != nil {
			h.OnError(err)
		}
		return fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 5. 检查状态码
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Webhook 返回错误: %d", resp.StatusCode)
		if h.OnError != nil {
			h.OnError(err)
		}
		return err
	}

	// 6. 检测响应类型并处理
	contentType := resp.Header.Get("Content-Type")

	if strings.Contains(contentType, "text/event-stream") {
		// SSE 格式（大厂标准）
		return h.readSSEStream(resp.Body)
	} else if resp.Header.Get("Transfer-Encoding") == "chunked" {
		// Chunked Transfer
		return h.readChunkedJSON(resp.Body)
	} else {
		// 普通 JSON（一次性响应）
		return h.readNormalJSON(resp.Body)
	}
}

// readSSEStream 读取 SSE 格式（支持多格式消息）
func (h *BotStreamHandler) readSSEStream(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer

	var currentData strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// 空行表示事件结束
		if line == "" {
			if currentData.Len() > 0 {
				// 解析片段
				chunk, err := parseBotChunk(currentData.String())
				if err != nil {
					return fmt.Errorf("解析片段失败: %w", err)
				}

				if h.OnChunk != nil {
					if err := h.OnChunk(chunk); err != nil {
						return err
					}
				}
				currentData.Reset()
			}
			continue
		}

		// 解析 data: 行
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			data = strings.TrimSpace(data)

			// 跳过 [DONE] 标记
			if data == "[DONE]" {
				if h.OnComplete != nil {
					return h.OnComplete()
				}
				return nil
			}

			currentData.WriteString(data)
		}
	}

	// 处理最后一个片段
	if currentData.Len() > 0 {
		chunk, err := parseBotChunk(currentData.String())
		if err != nil {
			return fmt.Errorf("解析片段失败: %w", err)
		}

		if h.OnChunk != nil {
			if err := h.OnChunk(chunk); err != nil {
				return err
			}
		}
	}

	// 完成回调
	if h.OnComplete != nil {
		return h.OnComplete()
	}

	return scanner.Err()
}

// readChunkedJSON 读取 JSON Lines 格式
func (h *BotStreamHandler) readChunkedJSON(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 解析 JSON
		chunk, err := parseBotChunk(line)
		if err != nil {
			// 如果不是 JSON，当作纯文本
			chunk = &BotChunk{
				Type:    "text",
				Content: line,
			}
		}

		if h.OnChunk != nil {
			if err := h.OnChunk(chunk); err != nil {
				return err
			}
		}
	}

	if h.OnComplete != nil {
		return h.OnComplete()
	}

	return scanner.Err()
}

// readNormalJSON 读取普通 JSON 响应
func (h *BotStreamHandler) readNormalJSON(reader io.Reader) error {
	var result map[string]interface{}
	if err := json.NewDecoder(reader).Decode(&result); err != nil {
		return fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 提取 content
	chunk := &BotChunk{
		Type:    "text",
		Content: "",
	}

	if content, ok := result["content"].(string); ok {
		chunk.Content = content
	} else if data, ok := result["data"].(map[string]interface{}); ok {
		if content, ok := data["content"].(string); ok {
			chunk.Content = content
		}
		if format, ok := data["format"].(string); ok {
			chunk.Type = format
		}
	}

	if h.OnChunk != nil {
		return h.OnChunk(chunk)
	}

	return nil
}

// parseBotChunk 解析 Bot 片段（支持多格式）
func parseBotChunk(data string) (*BotChunk, error) {
	// 尝试解析为结构化 JSON
	var chunk BotChunk
	if err := json.Unmarshal([]byte(data), &chunk); err == nil {
		// 默认类型为 text
		if chunk.Type == "" {
			chunk.Type = "text"
		}
		return &chunk, nil
	}

	// 如果不是 JSON，当作纯文本
	return &BotChunk{
		Type:    "text",
		Content: data,
	}, nil
}

// generateHMACSignature 生成 HMAC-SHA256 签名
func generateHMACSignature(body []byte, secret string, timestamp int64) string {
	// TODO: 实现真正的 HMAC-SHA256
	// 这里先用简单示例，生产环境应该用 crypto/hmac
	return fmt.Sprintf("sha256=%x", body)
}
