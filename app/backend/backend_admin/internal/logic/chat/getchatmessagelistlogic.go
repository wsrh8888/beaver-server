package logic

import (
	"context"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/user/user_models"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取聊天消息列表
func NewGetChatMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatMessageListLogic {
	return &GetChatMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatMessageListLogic) GetChatMessageList(req *types.GetChatMessageListReq) (resp *types.GetChatMessageListRes, err error) {
	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("1 = 1")

	// 会话ID筛选
	if req.ConversationID != "" {
		whereClause = whereClause.Where("conversation_id = ?", req.ConversationID)
	}

	// 发送者ID筛选
	if req.SendUserID != "" {
		whereClause = whereClause.Where("send_user_id = ?", req.SendUserID)
	}

	// 消息类型筛选
	if req.MsgType != 0 {
		whereClause = whereClause.Where("msg_type = ?", req.MsgType)
	}

	// 删除状态筛选 - GetChatMessageListReq 没有 Status 字段，使用 IsDeleted
	if req.IsDeleted {
		whereClause = whereClause.Where("status = ?", 4) // 4=已删除
	}

	// 时间范围筛选
	if req.StartTime != "" {
		if startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime); err == nil {
			whereClause = whereClause.Where("created_at >= ?", startTime)
		}
	}

	if req.EndTime != "" {
		if endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime); err == nil {
			whereClause = whereClause.Where("created_at <= ?", endTime)
		}
	}

	// 分页查询
	messages, count, err := list_query.ListQuery(l.svcCtx.DB, chat_models.ChatMessage{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.PageSize,
			Sort:  "created_at desc",
		},
		Where: whereClause,
	})

	if err != nil {
		logx.Errorf("查询聊天消息列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.GetChatMessageListItem
	for _, message := range messages {
		sendUserName := ""
		sendUserID := ""
		if message.SendUserID != nil && *message.SendUserID != "" {
			sendUserID = *message.SendUserID
			var user user_models.UserModel
			if err := l.svcCtx.DB.Where("uuid = ?", *message.SendUserID).First(&user).Error; err == nil {
				sendUserName = user.NickName
			}
		}

		list = append(list, types.GetChatMessageListItem{
			Id:             message.MessageID,
			MessageID:      message.MessageID,
			ConversationID: message.ConversationID,
			SendUserID:     sendUserID,
			SendUserName:   sendUserName,
			MsgType:        int(message.MsgType),
			MsgPreview:     message.MsgPreview,
			IsDeleted:      message.Status == 4, // 4=已删除
			CreateTime:     message.CreatedAt.String(),
			UpdateTime:     message.UpdatedAt.String(),
		})
	}

	return &types.GetChatMessageListRes{
		List:  list,
		Total: count,
	}, nil
}
