package logic

import (
	"context"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveBlacklistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveBlacklistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveBlacklistLogic {
	return &RemoveBlacklistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveBlacklistLogic) RemoveBlacklist(req *types.RemoveBlacklistReq) (resp *types.RemoveBlacklistRes, err error) {
	// todo: add your logic here and delete this line

	return
}
