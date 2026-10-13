/*
 * HTTP 工具：所有造数据脚本共用的 HTTP 客户端。
 *
 * 全局复用单个 http.Client，走连接池复用 TCP 连接，
 * 造大批量账号时避免反复握手；timeout 兜底防卡死。
 *
 * 独立成 package httputil（不叫 http，避免遮蔽 stdlib net/http）。
 */

package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     90 * time.Second,
	},
}

// PostJSON 发 POST JSON 请求，返回响应体与状态码。
// body 接受任意可序列化类型（map / struct），headers 为 nil 时只设默认 Content-Type。
func PostJSON(url string, body any, headers map[string]string) (string, int, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", 0, fmt.Errorf("序列化请求体: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return doReq(req)
}

// GetJSON 发 GET 请求，返回响应体与状态码。headers 设鉴权头等。
func GetJSON(url string, headers map[string]string) (string, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return doReq(req)
}

func doReq(req *http.Request) (string, int, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		return string(rb), resp.StatusCode, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(rb))
	}
	return string(rb), resp.StatusCode, nil
}
