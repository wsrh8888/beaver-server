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

type GetSyncMessageMediasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSyncMessageMediasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncMessageMediasLogic {
	return &GetSyncMessageMediasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncMessageMediasLogic) GetSyncMessageMedias(req *types.GetSyncMessageMediasReq) (*types.GetSyncMessageMediasRes, error) {
	var records []chat_models.ChatMessageMedia

	query := l.svcCtx.DB.Model(&chat_models.ChatMessageMedia{}).Where("user_id = ?", req.UserID)
	if req.Since > 0 {
		query = query.Where("UNIX_TIMESTAMP(created_at) > ?", req.Since)
	}

	if err := query.Find(&records).Error; err != nil {
		l.Logger.Errorf("同步消息媒体状态失败: userId=%s, error=%v", req.UserID, err)
		return nil, errors.New("同步失败")
	}

	messageIDs := make([]string, 0, len(records))
	for _, item := range records {
		messageIDs = append(messageIDs, item.MessageID)
	}

	return &types.GetSyncMessageMediasRes{
		MessageIDs:      messageIDs,
		ServerTimestamp: time.Now().Unix(),
	}, nil
}
