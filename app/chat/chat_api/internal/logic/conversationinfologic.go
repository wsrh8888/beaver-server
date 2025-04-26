package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/group/group_models"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ConversationInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConversationInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationInfoLogic {
	return &ConversationInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConversationInfoLogic) ConversationInfo(req *types.ConversationInfoReq) (resp *types.ConversationInfoRes, err error) {
	// 查询会话信息
	var userConversation chat_models.ChatUserConversationModel
	err = l.svcCtx.DB.Where("conversation_id = ? AND is_deleted = false", req.ConversationID).First(&userConversation).Error

	// 初始化响应
	resp = &types.ConversationInfoRes{
		ConversationID: req.ConversationID,
		MsgPreview:     "",
		UpdateAt:       "",
		IsTop:          false,
	}

	// 存在会话时填充消息预览和时间信息
	if err == nil {
		resp.MsgPreview = userConversation.LastMessage
		resp.UpdateAt = userConversation.UpdatedAt.String()
		resp.IsTop = userConversation.IsPinned
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 判断会话类型：包含下划线为私聊，否则为群聊
	if strings.Contains(req.ConversationID, "_") {
		// 私聊会话
		ids := strings.Split(req.ConversationID, "_")
		opponentID := ids[0]
		if ids[0] == req.UserID {
			opponentID = ids[1]
		}

		// 查询对方用户信息
		var user user_models.UserModel
		err = l.svcCtx.DB.Where("uuid = ?", opponentID).First(&user).Error
		if err != nil {
			return nil, err
		}

		resp.Avatar = user.Avatar
		resp.Nickname = user.NickName
		resp.ChatType = 1 // 私聊类型
	} else {
		// 群聊会话
		var group group_models.GroupModel
		err = l.svcCtx.DB.Where("uuid = ?", req.ConversationID).First(&group).Error
		if err != nil {
			return nil, err
		}

		resp.Avatar = group.Avatar
		resp.Nickname = group.Title
		resp.ChatType = 2 // 群聊类型
	}

	return resp, nil
}
