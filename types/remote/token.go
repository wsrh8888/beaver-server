// Package remote 对端接口的公共类型。
package remote

// TokenResponse 换票端点的回包（OAuth2 裸 JSON，不是 go-zero 信封）。
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
