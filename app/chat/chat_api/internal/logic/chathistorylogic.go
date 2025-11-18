package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/list_query"
	"beaver/common/models"
	"beaver/common/models/ctype"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatHistoryLogic {
	return &ChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatHistoryLogic) ChatHistory(req *types.ChatHistoryReq) (resp *types.ChatHistoryRes, err error) {

	fmt.Println("当前的会话Id是:", req.ConversationID)

	chatMessages, count, err := list_query.ListQuery(l.svcCtx.DB, chat_models.ChatMessage{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where: l.svcCtx.DB.Where("conversation_id = ?", req.ConversationID),
		// 移除Preload，因为微服务架构中不使用跨服务外键
	})

	if err != nil {
		return nil, err
	}

		// 收集需要查询用户信息的UserID列表（排除通知消息）
	var userIds []string
	userIdSet := make(map[string]bool)
	for _, chat := range chatMessages {
		if chat.SendUserID != nil && *chat.SendUserID != "" {
			if !userIdSet[*chat.SendUserID] {
				userIds = append(userIds, *chat.SendUserID)
				userIdSet[*chat.SendUserID] = true
			}
		}
	}

	// 批量获取用户信息
	userInfoMap := make(map[string]types.Sender)
	if len(userIds) > 0 {
		userListResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
			UserIdList: userIds,
		})
		if err != nil {
			l.Logger.Errorf("批量获取用户信息失败: %v", err)
			// 不返回错误，继续处理，为没有用户信息的消息设置默认值
		} else {
			// 转换用户信息
			for userId, userInfo := range userListResp.UserInfo {
				userInfoMap[userId] = types.Sender{
					UserID:   userId,
					Nickname: userInfo.NickName,
					Avatar:   userInfo.Avatar,
				}
			}
		}
	}

	// 构建消息历史
	var chatHistory []types.Message
	for _, chat := range chatMessages {
		// 转换消息内容
		var msg types.Msg
		if chat.Msg != nil {
			err := convertCtypeMsgToTypesMsg(*chat.Msg, &msg)
			if err != nil {
				return nil, err
			}
		}

		// 处理发送者信息
		var sender types.Sender
		sendUserID := ""
		if chat.SendUserID != nil {
			sendUserID = *chat.SendUserID
		}

		if sendUserID != "" {
			// 普通用户消息
			if userInfo, exists := userInfoMap[sendUserID]; exists {
				sender = userInfo
			} else {
				// 用户信息获取失败，使用默认值
				sender = types.Sender{
					UserID:   sendUserID,
					Nickname: "未知用户",
					Avatar:   "",
				}
			}
		} else {
			// 通知消息：SendUserID为空
			sender = types.Sender{
				UserID:   "",
				Nickname: "通知消息",
				Avatar:   "",
			}
		}

		message := types.Message{
			Id:               chat.Id,
			ConversationID:   chat.ConversationID,
			ConversationType: chat.ConversationType,
			Sender:           sender,
			CreateAt:         chat.CreatedAt.String(),
			Msg:              msg,
		}
		chatHistory = append(chatHistory, message)
	}

	return &types.ChatHistoryRes{
		Count: count,
		List:  chatHistory,
	}, nil
}

// Convert ctype.Msg to types.MsgType
func convertCtypeMsgToTypesMsg(ctypeMsg ctype.Msg, typesMsg *types.Msg) error {
	// 先用JSON解码为一个通用结构体
	encodedMsg, err := json.Marshal(ctypeMsg)
	if err != nil {
		return err
	}
	// 再将其解码为目标类型
	err = json.Unmarshal(encodedMsg, typesMsg)
	if err != nil {
		return err
	}
	return nil
}
