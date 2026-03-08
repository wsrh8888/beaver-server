package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/chat/chat_models"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncDeletedMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的已删除消息ID列表
func NewGetSyncDeletedMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncDeletedMessagesLogic {
	return &GetSyncDeletedMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncDeletedMessagesLogic) GetSyncDeletedMessages(req *types.GetSyncDeletedMessagesReq) (resp *types.GetSyncDeletedMessagesRes, err error) {
	var deletes []chat_models.ChatUserDelete

	// 1. 查询该用户在 Since 之后的删除记录
	query := l.svcCtx.DB.Model(&chat_models.ChatUserDelete{}).Where("user_id = ?", req.UserID)
	if req.Since > 0 {
		// 统一使用 Unix 时间戳进行同步比对
		query = query.Where("UNIX_TIMESTAMP(created_at) > ?", req.Since)
	}

	err = query.Find(&deletes).Error
	if err != nil {
		l.Logger.Errorf("同步已删除消息失败: userId=%s, error=%v", req.UserID, err)
		return nil, errors.New("同步失败")
	}

	// 2. 提取消息ID列表
	msgIDs := make([]string, 0, len(deletes))
	for _, d := range deletes {
		msgIDs = append(msgIDs, d.MessageID)
	}

	return &types.GetSyncDeletedMessagesRes{
		MessageIDs:      msgIDs,
		ServerTimestamp: time.Now().Unix(),
	}, nil
}
