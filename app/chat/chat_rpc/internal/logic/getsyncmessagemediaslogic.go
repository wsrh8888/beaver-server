package logic

import (
	"context"
	"time"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncMessageMediasLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSyncMessageMediasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncMessageMediasLogic {
	return &GetSyncMessageMediasLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSyncMessageMediasLogic) GetSyncMessageMedias(in *chat_rpc.GetSyncMessageMediasReq) (*chat_rpc.GetSyncMessageMediasRes, error) {
	var records []chat_models.ChatMessageMedia

	query := l.svcCtx.DB.Model(&chat_models.ChatMessageMedia{}).Where("user_id = ?", in.UserId)
	if in.Since > 0 {
		query = query.Where("UNIX_TIMESTAMP(created_at) > ?", in.Since)
	}

	if err := query.Find(&records).Error; err != nil {
		l.Errorf("同步消息媒体状态失败: userId=%s, error=%v", in.UserId, err)
		return nil, err
	}

	messageIDs := make([]string, 0, len(records))
	for _, item := range records {
		messageIDs = append(messageIDs, item.MessageID)
	}

	return &chat_rpc.GetSyncMessageMediasRes{
		MessageIds:      messageIDs,
		ServerTimestamp: time.Now().Unix(),
	}, nil
}
