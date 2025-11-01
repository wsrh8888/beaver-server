package logic

import (
	"context"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationVersionLogic {
	return &GetConversationVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationVersionLogic) GetConversationVersion(in *chat_rpc.GetConversationVersionReq) (*chat_rpc.GetConversationVersionRes, error) {
	var maxVersion int64

	// 查询该用户参与的所有会话的最大版本号
	// 通过关联用户会话设置表来找到用户参与的会话
	subQuery := l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Select("conversation_id").
		Where("user_id = ?", in.UserId)

	err := l.svcCtx.DB.Model(&chat_models.ChatConversationMeta{}).
		Where("conversation_id IN (?)", subQuery).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error

	if err != nil {
		l.Errorf("获取用户最新会话版本号失败: %v", err)
		return nil, err
	}

	return &chat_rpc.GetConversationVersionRes{
		LatestVersion: maxVersion,
	}, nil
}
