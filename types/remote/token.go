// Package remote 对端接口的公共类型。
package remote

// GetTokenApi 换票端点（跨 peer 共用协议结构，按规范例外单独成块，不套"每接口一块"）。
// 各 peer 的 remote/<peer>/token.go 统一用这两个类型解码 /machine/v1/token。
type (
	// TokenRequest 换票端点的请求体（client_credentials）。
	TokenRequest struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}

	// TokenResponse 换票端点的回包（标准 .api 信封 result）。
	TokenResponse struct {
		AccessToken string `json:"accessToken"`
		TokenType   string `json:"tokenType"`
		ExpiresIn   int    `json:"expiresIn"`
	}
)
