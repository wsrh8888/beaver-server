package logic

import (
	"context"
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateContentReportsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateContentReportsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateContentReportsLogic {
	return &UpdateContentReportsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateContentReportsLogic) UpdateContentReports(in *platform_rpc.UpdateContentReportsReq) (*platform_rpc.UpdateContentReportsRes, error) {
	if len(in.Ids) == 0 {
		return &platform_rpc.UpdateContentReportsRes{}, nil
	}

	updates := map[string]interface{}{"handler_id": in.HandlerId, "handle_remark": in.HandleRemark}
	switch in.Action {
	case 1:
		updates["status"] = platform_models.ReportStatusAccepted
		if in.CaseId > 0 {
			updates["case_id"] = in.CaseId
		}
	case 2:
		updates["status"] = platform_models.ReportStatusRejected
	case 3:
		updates["status"] = platform_models.ReportStatusResolved
	default:
		return nil, status.Error(codes.InvalidArgument, "无效的操作类型")
	}

	result := l.svcCtx.DB.Model(&platform_models.ContentReportModel{}).Where("id IN ?", in.Ids).Updates(updates)
	if result.Error != nil {
		l.Errorf("更新举报失败: %v", result.Error)
		return nil, status.Error(codes.Internal, "更新失败")
	}
	_ = time.Now()
	return &platform_rpc.UpdateContentReportsRes{AffectedCount: result.RowsAffected}, nil
}
