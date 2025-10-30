package logic

import (
	"context"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageSeqLogic {
	return &GetMessageSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessageSeqLogic) GetMessageSeq(in *chat_rpc.GetLatestSeqReq) (*chat_rpc.GetLatestSeqRes, error) {
	var maxSeq int64
	query := l.svcCtx.DB.Model(&chat_models.ChatMessage{}).Select("COALESCE(MAX(seq), 0)")

	if in.ConversationId != "" {
		query = query.Where("conversation_id = ?", in.ConversationId)
	}

	err := query.Scan(&maxSeq).Error
	if err != nil {
		l.Errorf("获取最新聊天序列号失败: %v", err)
		return nil, err
	}

	return &chat_rpc.GetLatestSeqRes{
		LatestSeq: maxSeq,
	}, nil
}
