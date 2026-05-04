package bot

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

	"beaver/app/open/open_models"
)

// StreamWebhookResponse 流式 Webhook 响应处理器
type StreamWebhookResponse struct {
	Config     *open_models.OpenWebhookConfig
	Payload    map[string]interface{}
	OnChunk    func(chunk string) error // 每收到一个片段的回调
	OnComplete func() error             // 完成时的回调
	OnError    func(err error) error    // 错误时的回调
}

// ExecuteStreamWebhook 执行流式 Webhook 调用
func ExecuteStreamWebhook(ctx context.Context, req *StreamWebhookResponse) error {
	// 1. 序列化 payload
	body, err := json.Marshal(req.Payload)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	// 2. 创建 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", req.Config.TargetURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 3. 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	// 4. 添加签名（如果有）
	if req.Config.Secret != "" {
		signature := generateSignature(body, req.Config.Secret)
		httpReq.Header.Set("X-Webhook-Signature", signature)
	}

	// 5. 发送请求
	client := &http.Client{
		Timeout: time.Duration(req.Config.Timeout) * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		if req.OnError != nil {
			req.OnError(err)
		}
		return fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 6. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Webhook 返回错误状态码: %d", resp.StatusCode)
		if req.OnError != nil {
			req.OnError(err)
		}
		return err
	}

	// 7. 读取流式响应
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		// SSE 格式
		return readSSEStream(resp.Body, req)
	} else {
		// 普通 JSON 格式，逐行读取
		return readJSONStream(resp.Body, req)
	}
}

// readSSEStream 读取 SSE 格式的流
func readSSEStream(reader io.Reader, req *StreamWebhookResponse) error {
	scanner := bufio.NewScanner(reader)
	var currentData strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// 空行表示事件结束
		if line == "" {
			if currentData.Len() > 0 {
				// 处理数据
				if req.OnChunk != nil {
					if err := req.OnChunk(currentData.String()); err != nil {
						return err
					}
				}
				currentData.Reset()
			}
			continue
		}

		// 解析 data: 行
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			currentData.WriteString(data)
			continue
		}
	}

	// 处理最后一个事件
	if currentData.Len() > 0 && req.OnChunk != nil {
		if err := req.OnChunk(currentData.String()); err != nil {
			return err
		}
	}

	// 完成回调
	if req.OnComplete != nil {
		return req.OnComplete()
	}

	return scanner.Err()
}

// readJSONStream 读取 JSON Lines 格式的流
func readJSONStream(reader io.Reader, req *StreamWebhookResponse) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// 尝试解析 JSON
		var chunkData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &chunkData); err != nil {
			// 如果不是 JSON，直接作为文本片段
			if req.OnChunk != nil {
				if err := req.OnChunk(line); err != nil {
					return err
				}
			}
			continue
		}

		// 提取 content 字段
		if content, ok := chunkData["content"].(string); ok {
			if req.OnChunk != nil {
				if err := req.OnChunk(content); err != nil {
					return err
				}
			}
		}
	}

	// 完成回调
	if req.OnComplete != nil {
		return req.OnComplete()
	}

	return scanner.Err()
}

// generateSignature 生成签名
func generateSignature(body []byte, secret string) string {
	// 简单的 HMAC 签名，实际生产环境应该用更安全的算法
	return fmt.Sprintf("%x", body)
}
