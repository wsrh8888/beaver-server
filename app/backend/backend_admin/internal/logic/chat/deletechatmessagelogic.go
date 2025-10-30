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

type DeleteChatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除聊天消息
func NewDeleteChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteChatMessageLogic {
	return &DeleteChatMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteChatMessageLogic) DeleteChatMessage(req *types.DeleteChatMessageReq) (resp *types.DeleteChatMessageRes, err error) {
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

	// 逻辑删除，设置is_deleted为true
	err = l.svcCtx.DB.Model(&message).Update("is_deleted", true).Error
	if err != nil {
		logx.Errorf("删除聊天消息失败: %v", err)
		return nil, errors.New("删除聊天消息失败")
	}

	return &types.DeleteChatMessageRes{}, nil
}
