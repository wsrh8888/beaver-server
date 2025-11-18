package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RestoreChatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 恢复已删除的消息
func NewRestoreChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreChatMessageLogic {
	return &RestoreChatMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestoreChatMessageLogic) RestoreChatMessage(req *types.RestoreChatMessageReq) (resp *types.RestoreChatMessageRes, err error) {
	// 检查消息是否存在
	var message chat_models.ChatMessage
	err = l.svcCtx.DB.Where("message_id = ?", req.MessageID).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("聊天消息不存在: %s", req.MessageID)
			return nil, errors.New("聊天消息不存在")
		}
		logx.Errorf("查询聊天消息失败: %v", err)
		return nil, errors.New("查询聊天消息失败")
	}

	// 恢复消息，设置is_deleted为false
	err = l.svcCtx.DB.Model(&message).Update("is_deleted", false).Error
	if err != nil {
		logx.Errorf("恢复聊天消息失败: %v", err)
		return nil, errors.New("恢复聊天消息失败")
	}

	return &types.RestoreChatMessageRes{}, nil
}
