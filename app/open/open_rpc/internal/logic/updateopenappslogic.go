package logic

import (
	"context"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateOpenAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOpenAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOpenAppsLogic {
	return &UpdateOpenAppsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateOpenAppsLogic) UpdateOpenApps(in *open_rpc.UpdateOpenAppsReq) (*open_rpc.UpdateOpenAppsRes, error) {
	if len(in.AppIds) == 0 {
		return &open_rpc.UpdateOpenAppsRes{}, nil
	}

	now := time.Now()
	updates := map[string]interface{}{"updated_at": now, "last_modified_by": in.OperatorId}

	switch in.Action {
	case 1: // 审核通过
		updates["audit_status"] = 1
		updates["status"] = 1
		updates["audited_by"] = in.OperatorId
		updates["audited_at"] = now
	case 2: // 审核拒绝
		updates["audit_status"] = 2
		updates["audited_by"] = in.OperatorId
		updates["audited_at"] = now
	case 3: // 禁用
		updates["status"] = 2
	case 4: // 启用（已发布）
		updates["status"] = 1
	default:
		return nil, status.Error(codes.InvalidArgument, "无效的操作类型")
	}

	result := l.svcCtx.DB.Model(&open_models.OpenApp{}).Where("app_id IN ?", in.AppIds).Updates(updates)
	if result.Error != nil {
		l.Errorf("更新应用状态失败: %v", result.Error)
		return nil, status.Error(codes.Internal, "操作失败")
	}

	return &open_rpc.UpdateOpenAppsRes{AffectedCount: result.RowsAffected}, nil
}
