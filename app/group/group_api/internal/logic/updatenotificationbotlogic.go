package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNotificationBotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateNotificationBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotificationBotLogic {
	return &UpdateNotificationBotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNotificationBotLogic) UpdateNotificationBot(req *types.UpdateNotificationBotReq) (resp *types.UpdateNotificationBotRes, err error) {
	var record open_models.OpenIncomingWebhook
	if err = l.svcCtx.DB.Where("id = ? AND app_id = ?", req.ID, "GROUP_NOTIFICATION").First(&record).Error; err != nil {
		return nil, errors.New("通知机器人不存在")
	}

	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", record.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可更新通知机器人")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Status != nil {
		if *req.Status != 0 && *req.Status != 1 {
			return nil, errors.New("状态值无效")
		}
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		return &types.UpdateNotificationBotRes{Success: true}, nil
	}

	if err = l.svcCtx.DB.Model(&record).Updates(updates).Error; err != nil {
		return nil, errors.New("更新失败")
	}

	return &types.UpdateNotificationBotRes{Success: true}, nil
}
