package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新机器人（名称/简介/头像/启用状态）
func NewUpdateBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBotLogic {
	return &UpdateBotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBotLogic) UpdateBot(req *types.UpdateBotReq) (resp *types.UpdateBotRes, err error) {
	// 1. 校验权限
	var ref group_models.GroupBotModel
	if err = l.svcCtx.DB.Where("bot_id = ?", req.BotID).First(&ref).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", ref.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可更新机器人")
	}

	// 2. 更新群内展示信息
	updates := make(map[string]interface{})
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) > 0 {
		if err := l.svcCtx.DB.Model(&ref).Updates(updates).Error; err != nil {
			return nil, errors.New("更新机器人信息失败")
		}
	}

	return &types.UpdateBotRes{}, nil
}
