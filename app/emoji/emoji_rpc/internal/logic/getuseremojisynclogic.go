package logic

import (
	"context"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserEmojiSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserEmojiSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserEmojiSyncLogic {
	return &GetUserEmojiSyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户表情数据同步接口（为微服务架构设计）
func (l *GetUserEmojiSyncLogic) GetUserEmojiSync(in *emoji_rpc.GetUserEmojiSyncReq) (*emoji_rpc.GetUserEmojiSyncRes, error) {
	// todo: add your logic here and delete this line

	return &emoji_rpc.GetUserEmojiSyncRes{}, nil
}
