// Package remote 对端接口的公共类型。
package remote

// TokenRequest 换票端点的请求体（client_credentials）。
type TokenRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// TokenResponse 换票端点的回包（标准 .api 信封 result）。
type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}
