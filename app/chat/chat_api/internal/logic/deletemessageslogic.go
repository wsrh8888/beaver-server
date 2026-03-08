package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除消息(仅自己不可见)
func NewDeleteMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessagesLogic {
	return &DeleteMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMessagesLogic) DeleteMessages(req *types.DeleteMessagesReq) (resp *types.DeleteMessagesRes, err error) {
	// 1. 批量记录删除记录
	// 对标大厂：这里不改主表状态，而是插入标记表，让同步和历史记录接口过滤
	deleteRecords := make([]chat_models.ChatUserDelete, 0)
	for _, msgID := range req.MessageIDs {
		// 生成版本号 (对标大厂：每条操作都有独立版本，确保多端增量同步的原子性)
		version := l.svcCtx.VersionGen.GetNextVersion("chat_user_deletes", "user_id", req.UserID)
		deleteRecords = append(deleteRecords, chat_models.ChatUserDelete{
			UserID:    req.UserID,
			MessageID: msgID,
			Version:   version,
		})
	}

	// 3. 执行入库 (使用 Create 批量插入)
	err = l.svcCtx.DB.Create(&deleteRecords).Error
	if err != nil {
		l.Logger.Errorf("用户 %s 批量删除消息失败: %v", req.UserID, err)
		return nil, err
	}

	return &types.DeleteMessagesRes{}, nil
}
