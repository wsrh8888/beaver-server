package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessagesByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按消息类型删除
func NewDeleteMessagesByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessagesByTypeLogic {
	return &DeleteMessagesByTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMessagesByTypeLogic) DeleteMessagesByType(req *types.DeleteMessagesByTypeReq) (resp *types.DeleteMessagesByTypeRes, err error) {
	if req.MsgType == 0 {
		return nil, errors.New("消息类型不能为空")
	}

	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("msg_type = ?", req.MsgType)

	// 会话ID筛选
	if req.ConversationID != "" {
		whereClause = whereClause.Where("conversation_id = ?", req.ConversationID)
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

	// 先统计要删除的数量
	var count int64
	err = whereClause.Model(&chat_models.ChatMessage{}).Count(&count).Error
	if err != nil {
		logx.Errorf("统计消息数量失败: %v", err)
		return nil, errors.New("统计消息数量失败")
	}

	// 逻辑删除消息
	err = whereClause.Update("is_deleted", true).Error
	if err != nil {
		logx.Errorf("按类型删除消息失败: %v", err)
		return nil, errors.New("按类型删除消息失败")
	}

	return &types.DeleteMessagesByTypeRes{
		DeletedCount: count,
	}, nil
}
