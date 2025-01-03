package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
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

	fmt.Println("当前的会话Id是:", req.ConversationId)

	chatMessages, count, err := list_query.ListQuery(l.svcCtx.DB, chat_models.ChatModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where:   l.svcCtx.DB.Where("conversation_id = ?", req.ConversationId),
		Preload: []string{"SendUserModel"},
	})

	if err != nil {
		return nil, err
	}

	// {"type":1,"textMsg":{"content":"你喜欢看体育比赛吗？"}}
	var chatHistory []types.Message
	for _, chat := range chatMessages {
		// 假设 ctype.Msg 和 types.MsgType 结构相同且纯粹需要类型转换
		var msg types.Msg
		err := convertCtypeMsgToTypesMsg(*chat.Msg, &msg)
		if err != nil {
			return nil, err
		}

		message := types.Message{
			MessageId:      chat.Id,
			ConversationId: chat.ConversationId,
			Sender: types.Sender{
				UserId:   chat.SendUserModel.UserId,
				Nickname: chat.SendUserModel.NickName,
				Avatar:   chat.SendUserModel.Avatar,
			},
			CreateAt: chat.CreatedAt.String(),
			Msg:      msg,
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
