/*
 * 公共 HTTP 工具：所有造数据脚本共用
 */

package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PostJSON 发 POST JSON 请求，返回响应体字符串
func PostJSON(url string, body map[string]string, headers map[string]string) (string, int, error) {
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return string(rb), resp.StatusCode, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(rb))
	}
	return string(rb), resp.StatusCode, nil
}
