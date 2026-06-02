package auth_public

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenRes, err error) {
	if req.RefreshToken == "" {
		return nil, errors.New("refreshToken 不能为空")
	}

	var oldToken open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("refresh_token = ?", req.RefreshToken).First(&oldToken).Error; err != nil {
		return nil, errors.New("刷新令牌无效")
	}
	if time.Now().Unix() > oldToken.RefreshTokenExpiresAt {
		return nil, errors.New("刷新令牌已过期，请重新授权")
	}

	rpcResp, err := l.svcCtx.OpenRpc.RefreshToken(l.ctx, &open_rpc.RefreshTokenReq{
		AppId:        oldToken.AppID,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		logx.Errorf("RefreshToken RPC 调用失败: %v", err)
		return nil, err
	}

	return &types.RefreshTokenRes{
		AccessToken:  rpcResp.AccessToken,
		RefreshToken: rpcResp.RefreshToken,
		ExpiresIn:    rpcResp.ExpiresIn,
	}, nil
}
