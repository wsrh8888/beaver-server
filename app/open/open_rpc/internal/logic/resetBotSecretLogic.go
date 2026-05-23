package logic

import (
	"context"

	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetBotSecretLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetBotSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetBotSecretLogic {
	return &ResetBotSecretLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetBotSecretLogic) ResetBotSecret(in *open_rpc.ResetBotSecretReq) (*open_rpc.ResetBotSecretRes, error) {
	// todo: add your logic here and delete this line

	return &open_rpc.ResetBotSecretRes{}, nil
}
