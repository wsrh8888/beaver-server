package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppLogic {
	return &AddAppLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AddAppLogic) AddApp(req *types.AddAppReq) (resp *types.AddAppRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.CreateApp(l.ctx, &platform_rpc.CreateAppReq{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		l.Errorf("??????: %v", err)
		return nil, err
	}
	return &types.AddAppRes{Id: uint(rpcRes.Id), AppID: rpcRes.AppId}, nil
}
