package moderation

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/backend/backend_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type EscalateContentReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEscalateContentReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EscalateContentReportLogic {
	return &EscalateContentReportLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *EscalateContentReportLogic) EscalateContentReport(req *types.EscalateContentReportReq) (resp *types.EscalateContentReportRes, err error) {
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
		return nil, errors.New("仅待处理举报可立案")
	}
	if report.CaseId > 0 {
		return nil, errors.New("该举报已关联工单")
	}

	priority := req.Priority
	if priority <= 0 {
		priority = 1
	}

	title := fmt.Sprintf("内容举报-%s", report.TargetId)
	caseRecord := backend_models.AdminModerationCase{
		CaseNo:      newCaseNo(),
		Source:      backend_models.CaseSourceReport,
		SourceID:    report.Id,
		TargetType:  int(report.TargetType),
		TargetID:    report.TargetId,
		Title:       title,
		Description: report.Content,
		Priority:    priority,
		Status:      backend_models.CaseStatusPending,
	}
	if err = l.svcCtx.DB.Create(&caseRecord).Error; err != nil {
		l.Errorf("创建工单失败: %v", err)
		return nil, err
	}

	_, err = l.svcCtx.PlatformRpc.UpdateContentReports(l.ctx, &platform_rpc.UpdateContentReportsReq{
		Ids:          []uint64{report.Id},
		Action:       1,
		CaseId:       uint64(caseRecord.Id),
		HandlerId:    req.UserID,
		HandleRemark: "举报立案",
	})
	if err != nil {
		l.Errorf("更新举报状态失败: %v", err)
		return nil, err
	}

	l.svcCtx.RecordOperation(req.UserID, "escalate_report", "report", fmt.Sprintf("%d", report.Id), uint64(caseRecord.Id),
		fmt.Sprintf("举报立案 targetType=%d targetId=%s", report.TargetType, report.TargetId), "success", "")

	return &types.EscalateContentReportRes{CaseID: uint64(caseRecord.Id), CaseNo: caseRecord.CaseNo}, nil
}
