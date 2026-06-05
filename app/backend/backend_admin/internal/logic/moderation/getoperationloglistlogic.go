package moderation

import (
	"context"
	"strings"

	"beaver/app/backend/backend_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperationLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperationLogListLogic {
	return &GetOperationLogListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetOperationLogListLogic) GetOperationLogList(req *types.GetOperationLogListReq) (resp *types.GetOperationLogListRes, err error) {
	page, pageSize := req.Page, req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	db := l.svcCtx.DB.Model(&backend_models.AdminOperationLog{})
	if req.OperatorID != "" {
		db = db.Where("operator_id = ?", req.OperatorID)
	}
	if req.TargetID != "" {
		db = db.Where("target_id = ?", req.TargetID)
	}
	if req.TargetType != "" {
		db = db.Where("target_type = ?", req.TargetType)
	}
	if req.Actions != "" {
		parts := strings.Split(req.Actions, ",")
		db = db.Where("action IN ?", parts)
	} else if req.Action != "" {
		db = db.Where("action = ?", req.Action)
	}
	if req.CaseID > 0 {
		db = db.Where("case_id = ?", req.CaseID)
	}

	var total int64
	if err = db.Count(&total).Error; err != nil {
		l.Errorf("统计审计日志失败: %v", err)
		return nil, err
	}

	var rows []backend_models.AdminOperationLog
	if err = db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		l.Errorf("查询审计日志失败: %v", err)
		return nil, err
	}

	list := make([]types.OperationLogInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, types.OperationLogInfo{
			ID:           uint64(row.Id),
			OperatorID:   row.OperatorID,
			Action:       row.Action,
			TargetType:   row.TargetType,
			TargetID:     row.TargetID,
			CaseID:       row.CaseID,
			Detail:       row.Detail,
			Result:       row.Result,
			ErrorMessage: row.ErrorMessage,
			CreatedAt:    row.CreatedAt.String(),
		})
	}
	return &types.GetOperationLogListRes{List: list, Total: total}, nil
}
