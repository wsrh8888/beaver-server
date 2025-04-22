package logic

import (
	"context"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBlacklistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBlacklistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlacklistLogic {
	return &GetBlacklistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBlacklistLogic) GetBlacklist(req *types.GetBlacklistReq) (resp *types.GetBlacklistRes, err error) {
	// todo: add your logic here and delete this line

	return
}
