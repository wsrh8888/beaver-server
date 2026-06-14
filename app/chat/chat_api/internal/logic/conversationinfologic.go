package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

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
	// 查询会话信息（这里应该是查询用户特定的会话设置）
	var userConversation chat_models.ChatUserConversation
	err = l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", req.ConversationID, req.UserID).First(&userConversation).Error

	// 初始化响应
	resp = &types.ConversationInfoRes{
		ConversationID: req.ConversationID,
		MsgPreview:     "",
		UpdatedAt:      "",
		IsTop:          false,
	}

	// 查询会话元数据获取最后消息
	var conversationMeta chat_models.ChatConversationMeta
	metaErr := l.svcCtx.DB.Where("conversation_id = ?", req.ConversationID).First(&conversationMeta).Error

	// 存在会话时填充消息预览和时间信息
	if err == nil {
		// 从会话元数据中获取最后消息
		if metaErr == nil {
			resp.MsgPreview = conversationMeta.LastMessage
		}
		resp.UpdatedAt = userConversation.UpdatedAt.String()
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

		userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
			UserID: opponentID,
		})
		if err != nil {
			return nil, err
		}
		if userRes.UserInfo == nil {
			return nil, errors.New("user not found")
		}

		resp.Avatar = userRes.UserInfo.Avatar
		resp.NickName = userRes.UserInfo.NickName
		resp.ChatType = 1 // 私聊类型
	} else {
		// 群聊会话
		groupRes, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{
			GroupIDs: []string{req.ConversationID},
		})
		if err != nil {
			return nil, err
		}
		if len(groupRes.Groups) == 0 {
			return nil, errors.New("group not found")
		}

		group := groupRes.Groups[0]
		resp.Avatar = group.Avatar
		resp.NickName = group.Name
		resp.ChatType = 2 // 群聊类型
	}

	return resp, nil
}
