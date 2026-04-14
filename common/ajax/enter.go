package ajax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 内部通讯密钥中心
const (
	// InternalSecret 内部服务互访密钥（如：Friend -> WS）
	InternalSecret = "beaver-internal-shared-secret-2026"
	// ExternalSecret 外部业务安全密钥
	ExternalSecret = "beaver-external-security-key-2026"
)

// WsInternalSecret 兼容旧逻辑，实际优先使用常量
var WsInternalSecret = InternalSecret

type ForwardRequest struct {
	ApiEndpoint string
	Method      string
	Token       string
	UserID      string
	Body        *bytes.Buffer
}

type Response struct {
	Code   int             `json:"code"`
	Msg    string          `json:"msg"`
	Result json.RawMessage `json:"result"`
}

func ForwardMessage(forwardReq ForwardRequest) (json.RawMessage, error) {
	client := &http.Client{}

	var req *http.Request
	var err error

	// 根据请求方法生成对应的HTTP请求
	if forwardReq.Method == "GET" {
		req, err = http.NewRequest("GET", forwardReq.ApiEndpoint, nil)
	} else if forwardReq.Method == "POST" {
		req, err = http.NewRequest("POST", forwardReq.ApiEndpoint, forwardReq.Body)
	} else {
		return nil, fmt.Errorf("不支持的请求方法: %s", forwardReq.Method)
	}

	if err != nil {
		return nil, fmt.Errorf("API请求创建错误: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", forwardReq.Token)
	req.Header.Set("Beaver-User-Id", forwardReq.UserID)
	req.Header.Set("X-Internal-Secret", WsInternalSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API请求错误: %v", err)
	}
	defer resp.Body.Close()

	// 检查API响应
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("消息转发未成功: %v", resp.Status)
	}

	// 读取API响应
	byteData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("消息转发错误: %v", err)
	}

	// 解析API响应
	var authResponse Response
	authErr := json.Unmarshal(byteData, &authResponse)
	if authErr != nil {
		return nil, fmt.Errorf("消息转发错误: %v", authErr)
	}

	if authResponse.Code != 0 {
		return nil, fmt.Errorf("消息转发失败: %v", authResponse.Msg)
	}
	return authResponse.Result, nil
}
