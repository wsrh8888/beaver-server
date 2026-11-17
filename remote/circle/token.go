package circle

import (
	"context"
	"errors"
	"fmt"

	"beaver/commonconfig"
	remotetype "beaver/types/remote"
	"beaver/utils/ajax"
)

// GetToken 用 client_credentials 换一次出站令牌，返回 accessToken 与有效期（秒，0 表示对端没给）。
//
// 换票是 beaver-social / social 的独立机器视图接口，地址与请求体都在这里，调用方只管拿票。
func GetTokenApi(ctx context.Context, clientID, clientSecret string) (string, int, error) {
	out, err := ajax.Post[remotetype.TokenResponse](ctx, commonconfig.SocialUrl, "/beaver-social/api/social/machine/v1/token", nil, remotetype.TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return "", 0, fmt.Errorf("换票失败: %w", err)
	}
	if out.AccessToken == "" {
		return "", 0, errors.New("换票响应缺少 accessToken")
	}
	return out.AccessToken, out.ExpiresIn, nil
}
