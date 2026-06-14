package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppsLogic {
	return &GetAppsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAppsLogic) GetApps(req *types.GetAppsReq) (resp *types.GetAppsRes, err error) {
	rpcReq := &platform_rpc.ListAppsReq{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	if req.IsActive {
		active := true
		rpcReq.IsActive = &active
	}

	rpcRes, err := l.svcCtx.PlatformRpc.ListApps(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("获取应用列表失败: %v", err)
		return nil, err
	}

	apps := make([]types.GetAppsItem, 0, len(rpcRes.Apps))
	for _, app := range rpcRes.Apps {
		apps = append(apps, types.GetAppsItem{
			Id:          uint(app.Id),
			AppID:       app.AppId,
			Name:        app.Name,
			Description: app.Description,
			IsActive:    app.IsActive,
			CreatedAt:   app.CreatedAt,
			UpdatedAt:   app.UpdatedAt,
		})
	}

	return &types.GetAppsRes{Total: rpcRes.Total, Apps: apps}, nil
}
