package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetNotificationBotSecretLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置通知机器人的签名密钥（旧 Secret 立即失效）
func NewResetNotificationBotSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetNotificationBotSecretLogic {
	return &ResetNotificationBotSecretLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetNotificationBotSecretLogic) ResetNotificationBotSecret(req *types.ResetNotificationBotSecretReq) (resp *types.ResetNotificationBotSecretRes, err error) {
	// 1. 查询本地展示信息
	var ref group_models.GroupBotModel
	if err = l.svcCtx.DB.First(&ref, req.ID).Error; err != nil {
		return nil, errors.New("通知机器人不存在")
	}

	// 2. 校验权限
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", ref.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可重置密钥")
	}

	// 3. 调 open_rpc 重置密钥（open 是 secret 的 master）
	rpcRes, err := l.svcCtx.OpenRpc.ResetBotSecret(l.ctx, &open_rpc.ResetBotSecretReq{
		Id: uint32(ref.BotID),
	})
	if err != nil {
		return nil, errors.New("重置失败")
	}

	return &types.ResetNotificationBotSecretRes{Secret: rpcRes.SignatureSecret}, nil
}
