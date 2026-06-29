package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UpdateWorkbenchAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkbenchAppLogic {
	return &UpdateWorkbenchAppLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateWorkbenchAppLogic) UpdateWorkbenchApp(in *platform_rpc.UpdateWorkbenchAppReq) (*platform_rpc.UpdateWorkbenchAppRes, error) {
	if strings.TrimSpace(in.WorkbenchAppId) == "" {
		return nil, status.Error(codes.InvalidArgument, "应用 ID 不能为空")
	}

	var app platform_models.WorkbenchApp
	if err := l.svcCtx.DB.Where("workbench_app_id = ?", in.WorkbenchAppId).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "工作台应用不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{
		"last_modified_by": in.OperatorId,
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		updates["name"] = name
	}
	if in.Description != "" {
		updates["description"] = strings.TrimSpace(in.Description)
	}
	if in.Icon != "" {
		updates["icon"] = strings.TrimSpace(in.Icon)
	}
	if entryURL := strings.TrimSpace(in.EntryUrl); entryURL != "" {
		updates["entry_url"] = entryURL
	}
	if in.Category != "" {
		updates["category"] = strings.TrimSpace(in.Category)
	}
	if in.Remark != "" {
		updates["remark"] = strings.TrimSpace(in.Remark)
	}
	if in.Sort != nil {
		updates["sort"] = int(*in.Sort)
	}
	if in.Status != nil {
		statusVal := int8(*in.Status)
		if statusVal != 0 && statusVal != 1 {
			return nil, status.Error(codes.InvalidArgument, "状态不合法")
		}
		updates["status"] = statusVal
	}

	if err := l.svcCtx.DB.Model(&app).Updates(updates).Error; err != nil {
		l.Errorf("更新工作台应用失败: %v", err)
		return nil, status.Error(codes.Internal, "更新工作台应用失败")
	}

	return &platform_rpc.UpdateWorkbenchAppRes{}, nil
}
