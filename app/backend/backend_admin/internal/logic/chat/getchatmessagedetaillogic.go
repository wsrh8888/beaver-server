package logic

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_models"

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

	err = l.svcCtx.DB.Preload("SendUserModel").Where("message_id = ?", req.MessageID).First(&message).Error
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
	if message.SendUserModel.NickName != "" {
		sendUserName = message.SendUserModel.NickName
	}
	if message.SendUserModel.FileName != "" {
		sendUserFileName = message.SendUserModel.FileName
	}

	msgContent := ""
	if message.Msg != nil {
		if msgBytes, err := json.Marshal(message.Msg); err == nil {
			msgContent = string(msgBytes)
		}
	}

	return &types.GetChatMessageDetailRes{
		Id:               message.MessageID,
		MessageID:        message.MessageID,
		ConversationID:   message.ConversationID,
		SendUserID:       message.SendUserID,
		SendUserName:     sendUserName,
		SendUserFileName: sendUserFileName,
		MsgType:          int(message.MsgType),
		MsgPreview:       message.MsgPreview,
		MsgContent:       msgContent,
		IsDeleted:        message.IsDeleted,
		CreateTime:       message.CreatedAt.String(),
		UpdateTime:       message.UpdatedAt.String(),
	}, nil
}
