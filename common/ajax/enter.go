package ajax

import (
	"beaver/common/etcd"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForwardRequest struct {
	ApiEndpoint string
	Method      string
	Token       string
	UserId      string
	Body        *bytes.Buffer
}

type Response struct {
	Code   int             `json:"code"`
	Msg    string          `json:"msg"`
	Result json.RawMessage `json:"result"`
}

type WsProxyReq struct {
	UserId   string                 `header:"Beaver-User-Id"`
	Command  string                 `json:"command"`
	TargetId string                 `json:"targetId"`
	Type     string                 `json:"type"`
	Body     map[string]interface{} `json:"body"`
}

func SendMessageToWs(etcdUrl string, types string, senderId string, targetId string, requestBody map[string]interface{}) {
	addr := etcd.GetServiceAddr(etcdUrl, "ws_api")
	if addr == "" {
		logx.Error("未匹配到服务")
		return
	}
	apiEndpoint := fmt.Sprintf("http://%s/api/ws/proxySendMsg", addr)

	wsProxyReq := WsProxyReq{
		UserId:   senderId,
		Command:  "COMMON_UPDATE_MESSAGE",
		TargetId: targetId,
		Type:     types,
		Body:     requestBody,
	}
	body, _ := json.Marshal(wsProxyReq)

	ForwardMessage(ForwardRequest{
		ApiEndpoint: apiEndpoint,
		Method:      "POST",
		Token:       "",
		UserId:      senderId,
		Body:        bytes.NewBuffer(body),
	})
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
	req.Header.Set("Token", forwardReq.Token)           // 使用Token进行鉴权
	req.Header.Set("Beaver-User-Id", forwardReq.UserId) // 使用Token进行鉴权

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
	sendAjaxJSON, _ := json.Marshal(authResponse.Result)

	fmt.Println("消息转发成功", string(sendAjaxJSON))
	return authResponse.Result, nil
}
