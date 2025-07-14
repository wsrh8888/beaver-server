package ajax

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetCityByIP 根据IP获取城市名称
func GetCityByIP(ip string) (string, error) {
	// 调用IP-API.com免费服务
	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=en", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Status  string `json:"status"`
		City    string `json:"city"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Status != "success" {
		return "", fmt.Errorf("IP定位失败: %s", result.Message)
	}

	return result.City, nil
}
