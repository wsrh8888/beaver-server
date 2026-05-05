package app

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetAppSecretLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置应用密钥
func NewResetAppSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetAppSecretLogic {
	return &ResetAppSecretLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetAppSecretLogic) ResetAppSecret(req *types.ResetAppSecretReq) (resp *types.ResetAppSecretRes, err error) {
	// todo: add your logic here and delete this line

	return
}
