package oauth_secret

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRevokeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeTokenLogic {
	return &RevokeTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeTokenLogic) RevokeToken(req *types.RevokeTokenReq) (resp *types.RevokeTokenRes, err error) {
	if req.Token == "" {
		return nil, errors.New("token 不能为空")
	}

	result := l.svcCtx.DB.Where("token = ? OR refresh_token = ?", req.Token, req.Token).Delete(&open_models.OpenOAuthToken{})
	if result.Error != nil {
		return nil, errors.New("撤销令牌失败")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("令牌不存在")
	}

	return &types.RevokeTokenRes{Success: true}, nil
}
