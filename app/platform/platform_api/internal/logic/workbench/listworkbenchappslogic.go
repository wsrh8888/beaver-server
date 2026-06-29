package workbench

import (
	"context"
	"errors"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkbenchAppsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWorkbenchAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkbenchAppsLogic {
	return &ListWorkbenchAppsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWorkbenchAppsLogic) ListWorkbenchApps(req *types.ListWorkbenchAppsReq) (*types.ListWorkbenchAppsRes, error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	rpcRes, err := l.svcCtx.PlatformRpc.ListEnabledWorkbenchApps(l.ctx, &platform_rpc.ListEnabledWorkbenchAppsReq{
		Category: req.Category,
	})
	if err != nil {
		l.Errorf("获取工作台应用列表失败: %v", err)
		return nil, errors.New("获取工作台应用列表失败")
	}

	list := make([]types.ListWorkbenchAppsItem, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		list = append(list, types.ListWorkbenchAppsItem{
			WorkbenchAppID: item.WorkbenchAppId,
			Name:           item.Name,
			Description:    item.Description,
			Icon:           item.Icon,
			EntryURL:       item.EntryUrl,
			Category:       item.Category,
			Sort:           int(item.Sort),
		})
	}

	return &types.ListWorkbenchAppsRes{List: list}, nil
}
