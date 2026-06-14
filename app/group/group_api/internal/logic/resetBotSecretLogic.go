package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)


type ResetBotSecretLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 重置机器人的签名密钥（旧 Secret 立即失效）
func NewResetBotSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetBotSecretLogic {
	return &ResetBotSecretLogic{
		ctx:    ctx,
		logger: logger.New("reset_bot_secret"),
		svcCtx: svcCtx,
	}
}

func (l *ResetBotSecretLogic) ResetBotSecret(req *types.ResetBotSecretReq) (resp *types.ResetBotSecretRes, err error) {
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
		return nil, errors.New("无权限，仅群主或管理员可重置密钥")
	}

	// 2. 通过 open_rpc 获取 Bot ID
	botInfoRes, err := l.svcCtx.OpenRpc.GetBotInfo(l.ctx, &open_rpc.GetBotInfoReq{
		BotId: ref.BotID,
	})
	if err != nil {
		return nil, errors.New("Open Bot 记录不存在")
	}

	// 3. 调用 open_rpc 重置密钥
	secretRes, err := l.svcCtx.OpenRpc.ResetBotSecret(l.ctx, &open_rpc.ResetBotSecretReq{
		Id: botInfoRes.Id,
	})
	if err != nil {
		return nil, errors.New("重置密钥失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "群机器人密钥重置成功",
		Data: map[string]interface{}{
			"groupId": ref.GroupID,
			"userId":  req.UserID,
			"botId":   req.BotID,
		},
	})

	return &types.ResetBotSecretRes{
		Secret: secretRes.SignatureSecret,
	}, nil
}
