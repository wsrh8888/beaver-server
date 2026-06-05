package moderation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"beaver/app/backend/backend_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type HandleModerationCaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleModerationCaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleModerationCaseLogic {
	return &HandleModerationCaseLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *HandleModerationCaseLogic) HandleModerationCase(req *types.HandleModerationCaseReq) (resp *types.HandleModerationCaseRes, err error) {
	if req.CaseID == 0 {
		return nil, errors.New("工单ID不能为空")
	}
	if req.Status < backend_models.CaseStatusPending || req.Status > backend_models.CaseStatusRejected {
		return nil, errors.New("无效的工单状态")
	}

	var c backend_models.AdminModerationCase
	if err = l.svcCtx.DB.Where("id = ?", req.CaseID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工单不存在")
		}
		l.Errorf("查询工单失败: %v", err)
		return nil, err
	}

	for _, act := range req.Actions {
		if actErr := executeControlAction(l.ctx, l.svcCtx, req.UserID, uint64(c.Id), act); actErr != nil {
			l.Errorf("执行管控动作失败 action=%s: %v", act.Action, actErr)
			return nil, actErr
		}
	}

	actionsJSON := ""
	if len(req.Actions) > 0 {
		if b, mErr := json.Marshal(req.Actions); mErr == nil {
			actionsJSON = string(b)
		}
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":        req.Status,
		"handler_id":    req.UserID,
		"handle_remark": req.HandleRemark,
		"handle_time":   &now,
		"actions_taken": actionsJSON,
	}
	if err = l.svcCtx.DB.Model(&c).Updates(updates).Error; err != nil {
		l.Errorf("更新工单失败: %v", err)
		return nil, err
	}

	if req.Status == backend_models.CaseStatusResolved || req.Status == backend_models.CaseStatusRejected {
		reportRes, listErr := l.svcCtx.PlatformRpc.ListContentReports(l.ctx, &platform_rpc.ListContentReportsReq{
			Page:       1,
			PageSize:   100,
			TargetType: int32(c.TargetType),
			TargetId:   c.TargetID,
		})
		if listErr == nil && reportRes != nil {
			ids := make([]uint64, 0)
			for _, r := range reportRes.List {
				if r.CaseId == uint64(c.Id) || r.CaseId == 0 {
					ids = append(ids, r.Id)
				}
			}
			if len(ids) > 0 {
				action := int32(3)
				if req.Status == backend_models.CaseStatusRejected {
					action = 2
				}
				_, _ = l.svcCtx.PlatformRpc.UpdateContentReports(l.ctx, &platform_rpc.UpdateContentReportsReq{
					Ids:          ids,
					Action:       action,
					HandlerId:    req.UserID,
					HandleRemark: req.HandleRemark,
				})
			}
		}
	}

	l.svcCtx.RecordOperation(req.UserID, "handle_case", "case", fmt.Sprintf("%d", c.Id), uint64(c.Id),
		fmt.Sprintf("处置工单 status=%d remark=%s", req.Status, req.HandleRemark), "success", "")

	return &types.HandleModerationCaseRes{}, nil
}
