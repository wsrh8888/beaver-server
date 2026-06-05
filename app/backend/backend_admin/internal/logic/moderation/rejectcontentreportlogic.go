package moderation

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectContentReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectContentReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectContentReportLogic {
	return &RejectContentReportLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RejectContentReportLogic) RejectContentReport(req *types.RejectContentReportReq) (resp *types.RejectContentReportRes, err error) {
	if req.ReportID == 0 {
		return nil, errors.New("举报ID不能为空")
	}

	reportRes, err := l.svcCtx.PlatformRpc.GetContentReport(l.ctx, &platform_rpc.GetContentReportReq{Id: req.ReportID})
	if err != nil {
		l.Errorf("获取举报详情失败: %v", err)
		return nil, err
	}
	if reportRes.Report == nil {
		return nil, errors.New("举报不存在")
	}
	report := reportRes.Report
	if report.Status != platform_models.ReportStatusPending {
		return nil, errors.New("仅待处理举报可驳回")
	}

	remark := req.HandleRemark
	if remark == "" {
		remark = "举报驳回"
	}

	_, err = l.svcCtx.PlatformRpc.UpdateContentReports(l.ctx, &platform_rpc.UpdateContentReportsReq{
		Ids:          []uint64{report.Id},
		Action:       2,
		HandlerId:    req.UserID,
		HandleRemark: remark,
	})
	if err != nil {
		l.Errorf("驳回举报失败: %v", err)
		return nil, err
	}

	l.svcCtx.RecordOperation(req.UserID, "reject_report", "report", fmt.Sprintf("%d", report.Id), 0,
		fmt.Sprintf("驳回举报 targetType=%d targetId=%s remark=%s", report.TargetType, report.TargetId, remark), "success", "")

	return &types.RejectContentReportRes{}, nil
}
