package logic

import (
	"context"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateTokenLogic {
	return &ValidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ValidateTokenLogic) ValidateToken(in *open_rpc.ValidateTokenReq) (*open_rpc.ValidateTokenRes, error) {
	// 查询 Access Token
	var token open_models.OpenAccessToken
	if err := l.svcCtx.DB.Where("token = ?", in.AccessToken).First(&token).Error; err != nil {
		return &open_rpc.ValidateTokenRes{
			Valid: false,
		}, nil
	}

	// 检查 Token 是否过期
	valid := time.Now().Unix() <= token.ExpiresAt

	return &open_rpc.ValidateTokenRes{
		Valid:  valid,
		UserId: token.UserID,
		AppId:  token.AppID,
	}, nil
}
