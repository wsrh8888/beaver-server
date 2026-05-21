package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_models"

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
	var ref group_models.GroupNotificationBotModel
	if err = l.svcCtx.DB.First(&ref, req.ID).Error; err != nil {
		return nil, errors.New("通知机器人不存在")
	}

	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", ref.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可更新通知机器人")
	}

	// 调 open_rpc 更新 master 记录
	if req.Status != nil {
		if *req.Status != 0 && *req.Status != 1 {
			return nil, errors.New("状态值无效")
		}
		if _, err = l.svcCtx.OpenRpc.UpdateWebhook(l.ctx, &open_rpc.UpdateWebhookReq{
			Id:     uint32(ref.WebhookID),
			Status: int32(*req.Status),
		}); err != nil {
			return nil, errors.New("更新失败")
		}
	}

	// 同步本地引用表
	localUpdates := map[string]interface{}{}
	if req.Name != "" {
		localUpdates["name"] = req.Name
	}
	if req.Description != "" {
		localUpdates["description"] = req.Description
	}
	if req.Avatar != "" {
		localUpdates["avatar"] = req.Avatar
	}
	if req.Type != "" {
		localUpdates["type"] = req.Type
	}
	if req.Status != nil {
		localUpdates["status"] = *req.Status
	}
	if len(localUpdates) > 0 {
		l.svcCtx.DB.Model(&ref).Updates(localUpdates)
	}

	// 同步机器人用户的昵称和头像，保证聊天界面显示最新信息
	botUpdates := map[string]interface{}{}
	if req.Name != "" {
		botUpdates["nick_name"] = req.Name
	}
	if req.Avatar != "" {
		botUpdates["avatar"] = req.Avatar
	}
	if len(botUpdates) > 0 {
		l.svcCtx.DB.Model(&user_models.UserModel{}).Where("user_id = ?", ref.BotUserID).Updates(botUpdates)
	}

	return &types.UpdateNotificationBotRes{Success: true}, nil
}
