package handler

import (
	"encoding/json"
	"net/http"

	"beaver/app/datasync/datasync_api/internal/config"
	"beaver/middleware/jwt"
)

// TokenHandler 内部编排面的 /internal/token 端点：OAuth2 client_credentials 换发 JWT access token。
//
// 调用方（收费模块 beaver-social 等）拿 client_id + client_secret 来换；校验通过用本面 S2S.Secret
// 签发一个受众为 server:internal 的 Bearer。该端点本身不要求 Bearer（否则换不到），靠 client_secret 鉴权。
func TokenHandler(c config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		clientID, clientSecret := parseClientCredentials(r)
		scopes, ok := jwt.AuthenticateClient(c.S2S.Clients, clientID, clientSecret)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid_client"})
			return
		}
		signed, err := jwt.Issue(c.S2S.Secret, c.Name, clientID, "server:internal", scopes, jwt.S2STTL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": signed,
			"token_type":   "Bearer",
			"expires_in":   int(jwt.S2STTL.Seconds()),
		})
	}
}

// parseClientCredentials 从 Basic Auth 或表单取 client_id / client_secret。
func parseClientCredentials(r *http.Request) (string, string) {
	if id, secret, ok := r.BasicAuth(); ok {
		return id, secret
	}
	_ = r.ParseForm()
	return r.PostForm.Get("client_id"), r.PostForm.Get("client_secret")
}
