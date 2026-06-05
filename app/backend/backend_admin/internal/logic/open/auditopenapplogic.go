package open

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditOpenAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuditOpenAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditOpenAppLogic {
	return &AuditOpenAppLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AuditOpenAppLogic) AuditOpenApp(req *types.AuditOpenAppReq) (resp *types.AuditOpenAppRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.AppID == "" {
		return nil, errors.New("应用ID不能为空")
	}
	if req.Status != 1 && req.Status != 2 {
		return nil, errors.New("无效的审核状态")
	}

	_, err = l.svcCtx.OpenRpc.UpdateOpenApps(l.ctx, &open_rpc.UpdateOpenAppsReq{
		AppIds:     []string{req.AppID},
		Action:     int32(req.Status),
		OperatorId: req.UserID,
		AuditRemark: req.AuditRemark,
	})
	if err != nil {
		l.Errorf("审核应用失败: %v", err)
		return nil, err
	}

	return &types.AuditOpenAppRes{}, nil
}
