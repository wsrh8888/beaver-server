package circle

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"beaver/commonconfig"
	remotetype "beaver/types/remote"
	"beaver/utils/ajax"
)

// GetToken 用 client_credentials 换一次出站令牌，返回 access_token 与有效期（秒，0 表示对端没给）。
//
// 换票是收费模块圈子服务的一个独立接口，地址和表单都在这里，调用方只管拿票。
func GetToken(ctx context.Context, clientID, clientSecret string) (string, int, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	out, err := ajax.PostForm[remotetype.TokenResponse](ctx, commonconfig.CircleUrl+"/internal/token", nil, form)
	if err != nil {
		return "", 0, fmt.Errorf("换票失败: %w", err)
	}
	if out.AccessToken == "" {
		return "", 0, errors.New("换票响应缺少 access_token")
	}
	return out.AccessToken, out.ExpiresIn, nil
}
