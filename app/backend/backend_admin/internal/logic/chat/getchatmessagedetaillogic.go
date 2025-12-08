package logic

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetChatMessageDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取聊天消息详情
func NewGetChatMessageDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatMessageDetailLogic {
	return &GetChatMessageDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatMessageDetailLogic) GetChatMessageDetail(req *types.GetChatMessageDetailReq) (resp *types.GetChatMessageDetailRes, err error) {
	var message chat_models.ChatMessage

	err = l.svcCtx.DB.Where("message_id = ?", req.MessageID).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("聊天消息不存在: %s", req.MessageID)
			return nil, errors.New("聊天消息不存在")
		}
		logx.Errorf("查询聊天消息详情失败: %v", err)
		return nil, errors.New("查询聊天消息详情失败")
	}

	sendUserName := ""
	sendUserFileName := ""
	if message.SendUserID != nil && *message.SendUserID != "" {
		var user user_models.UserModel
		if err := l.svcCtx.DB.Where("user_id = ?", *message.SendUserID).First(&user).Error; err == nil {
			sendUserName = user.NickName
			sendUserFileName = user.Avatar
		}
	}

	msgContent := ""
	if message.Msg != nil {
		if msgBytes, err := json.Marshal(message.Msg); err == nil {
			msgContent = string(msgBytes)
		}
	}

	sendUserID := ""
	if message.SendUserID != nil {
		sendUserID = *message.SendUserID
	}

	return &types.GetChatMessageDetailRes{
		Id:               message.MessageID,
		MessageID:        message.MessageID,
		ConversationID:   message.ConversationID,
		SendUserID:       sendUserID,
		SendUserName:     sendUserName,
		SendUserFileName: sendUserFileName,
		MsgType:          int(message.MsgType),
		MsgPreview:       message.MsgPreview,
		MsgContent:       msgContent,
		IsDeleted:        message.Status == 4, // 4=已删除
		CreateTime:       message.CreatedAt.String(),
		UpdateTime:       message.UpdatedAt.String(),
	}, nil
}
