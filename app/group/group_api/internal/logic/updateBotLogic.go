package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

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

	if err = l.updateBotUserProfile(req); err != nil {
		return nil, err
	}

	if err = l.updateBotSecurity(req); err != nil {
		return nil, err
	}

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

func (l *UpdateBotLogic) updateBotUserProfile(req *types.UpdateBotReq) error {
	if req.Name == "" && req.Avatar == "" && req.Description == "" {
		return nil
	}

	patchReq := &user_rpc.UpdateUsersReq{
		UserIds: []string{req.BotID},
	}
	if req.Name != "" {
		patchReq.PatchNickName = &req.Name
	}
	if req.Avatar != "" {
		patchReq.PatchAvatar = &req.Avatar
	}
	if req.Description != "" {
		patchReq.PatchAbstract = &req.Description
	}

	_, err := l.svcCtx.UserRpc.UpdateUsers(l.ctx, patchReq)
	if err != nil {
		l.Logger.Errorf("更新机器人用户资料失败: botId=%s, error=%v", req.BotID, err)
		return errors.New("更新机器人资料失败")
	}
	return nil
}

func (l *UpdateBotLogic) updateBotSecurity(req *types.UpdateBotReq) error {
	if req.Security == nil {
		return nil
	}

	if req.Security.KeywordsEnabled && len(req.Security.Keywords) > 10 {
		return errors.New("关键词最多10个")
	}

	_, err := l.svcCtx.OpenRpc.UpdateBot(l.ctx, &open_rpc.UpdateBotReq{
		BotId: req.BotID,
		Security: &open_rpc.BotSecurity{
			KeywordsEnabled:    req.Security.KeywordsEnabled,
			Keywords:           req.Security.Keywords,
			IpWhitelistEnabled: req.Security.IPWhitelistEnabled,
			IpWhitelist:        req.Security.IPWhitelist,
			SignatureEnabled:   req.Security.SignatureEnabled,
			SignatureSecret:    req.Security.SignatureSecret,
		},
	})
	if err != nil {
		l.Logger.Errorf("更新机器人安全设置失败: botId=%s, error=%v", req.BotID, err)
		return errors.New("更新机器人安全设置失败")
	}
	return nil
}
