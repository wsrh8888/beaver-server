package moderation

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/backend/backend_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateModerationCaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateModerationCaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateModerationCaseLogic {
	return &CreateModerationCaseLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateModerationCaseLogic) CreateModerationCase(req *types.CreateModerationCaseReq) (resp *types.CreateModerationCaseRes, err error) {
	if req.TargetType <= 0 || req.TargetID == "" {
		return nil, errors.New("处置对象不能为空")
	}
	if req.Title == "" {
		return nil, errors.New("工单标题不能为空")
	}

	priority := req.Priority
	if priority <= 0 {
		priority = 1
	}

	caseRecord := backend_models.AdminModerationCase{
		CaseNo:      newCaseNo(),
		Source:      backend_models.CaseSourceManual,
		TargetType:  req.TargetType,
		TargetID:    req.TargetID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		Status:      backend_models.CaseStatusPending,
	}
	if err = l.svcCtx.DB.Create(&caseRecord).Error; err != nil {
		l.Errorf("创建工单失败: %v", err)
		return nil, err
	}

	l.svcCtx.RecordOperation(req.UserID, "create_case", "case", fmt.Sprintf("%d", caseRecord.Id), uint64(caseRecord.Id),
		fmt.Sprintf("手动创建工单 targetType=%d targetId=%s", req.TargetType, req.TargetID), "success", "")

	return &types.CreateModerationCaseRes{CaseID: uint64(caseRecord.Id), CaseNo: caseRecord.CaseNo}, nil
}
