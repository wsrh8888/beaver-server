package logic

import (
	"context"
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 聊天数据同步
func NewChatSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatSyncLogic {
	return &ChatSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatSyncLogic) ChatSync(req *types.ChatSyncReq) (resp *types.ChatSyncRes, err error) {
	var chats []chat_models.ChatMessage

	// 构建基础查询条件
	var query = l.svcCtx.DB.Where("seq > ? AND seq <= ?", req.FromSeq, req.ToSeq)

	// 如果指定了会话ID，则只同步该会话的消息
	if req.ConversationID != "" {
		query = query.Where("conversation_id = ?", req.ConversationID)
	} else {
		// 如果没有指定会话ID，需要过滤出用户相关的会话
		// 通过子查询获取用户参与的所有会话ID（包括已删除的）
		subQuery := l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
			Select("conversation_id").
			Where("user_id = ?", req.UserID)

		query = query.Where("conversation_id IN (?)", subQuery)
	}

	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	// 查询数据
	err = query.Order("seq ASC").Limit(limit + 1).Find(&chats).Error
	if err != nil {
		l.Errorf("查询聊天数据失败: %v", err)
		return nil, err
	}

	// 判断是否还有更多数据
	hasMore := len(chats) > limit
	if hasMore {
		chats = chats[:limit]
	}

	// 转换为响应格式
	var messages = make([]types.ChatSyncMessage, 0)
	var nextSeq int64 = req.FromSeq

	for _, chat := range chats {
		// 序列化消息内容
		msgJson := ""
		if chat.Msg != nil {
			// 这里需要将Msg结构体序列化为JSON字符串
			// 由于Msg是*ctype.Msg类型，需要根据实际类型进行序列化
			msgJson = chat.MsgPreview // 暂时使用预览作为消息内容
		}

		messages = append(messages, types.ChatSyncMessage{
			MessageID:      chat.MessageID,
			ConversationID: chat.ConversationID,
			SendUserID:     chat.SendUserID,
			MsgType:        uint32(chat.MsgType),
			MsgPreview:     chat.MsgPreview,
			Msg:            msgJson,
			IsDeleted:      chat.IsDeleted,
			Seq:            chat.Seq,
			CreateAt:       time.Time(chat.CreatedAt).Unix(),
		})

		nextSeq = chat.Seq
	}

	// 如果没有更多数据，nextSeq应该是toSeq+1
	if !hasMore {
		nextSeq = req.ToSeq + 1
	}

	resp = &types.ChatSyncRes{
		Messages: messages,
		HasMore:  hasMore,
		NextSeq:  nextSeq,
	}

	l.Infof("聊天数据同步完成，用户ID: %s, 返回消息数: %d, 还有更多: %v", req.UserID, len(messages), hasMore)
	return resp, nil
}
