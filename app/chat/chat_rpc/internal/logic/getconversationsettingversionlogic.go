package logic

import (
	"context"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationSettingVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationSettingVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationSettingVersionLogic {
	return &GetConversationSettingVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationSettingVersionLogic) GetConversationSettingVersion(in *chat_rpc.GetConversationSettingVersionReq) (*chat_rpc.GetConversationSettingVersionRes, error) {
	var maxVersion int64
	err := l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Where("user_id = ?", in.UserId).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error

	if err != nil {
		l.Errorf("获取最新会话版本号失败: %v", err)
		return nil, err
	}

	return &chat_rpc.GetConversationSettingVersionRes{
		LatestVersion: maxVersion,
	}, nil
}
