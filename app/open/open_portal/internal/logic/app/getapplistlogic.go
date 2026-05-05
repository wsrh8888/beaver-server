package app

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取应用列表
func NewGetAppListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppListLogic {
	return &GetAppListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppListLogic) GetAppList(req *types.GetAppListReq) (resp *types.GetAppListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
