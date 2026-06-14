package open

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOpenAppStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOpenAppStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOpenAppStatusLogic {
	return &UpdateOpenAppStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateOpenAppStatusLogic) UpdateOpenAppStatus(req *types.UpdateOpenAppStatusReq) (resp *types.UpdateOpenAppStatusRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if len(req.AppIDs) == 0 {
		return nil, errors.New("请选择应用")
	}
	if req.Action != 3 && req.Action != 4 {
		return nil, errors.New("无效的操作类型")
	}

	_, err = l.svcCtx.OpenRpc.UpdateOpenApps(l.ctx, &open_rpc.UpdateOpenAppsReq{
		AppIds:     req.AppIDs,
		Action:     int32(req.Action),
		OperatorId: req.UserID,
	})
	if err != nil {
		l.Errorf("更新应用状态失败: %v", err)
		return nil, err
	}

	return &types.UpdateOpenAppStatusRes{}, nil
}
