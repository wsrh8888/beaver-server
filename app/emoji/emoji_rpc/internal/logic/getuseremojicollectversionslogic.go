package logic

import (
	"context"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserEmojiCollectVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserEmojiCollectVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserEmojiCollectVersionsLogic {
	return &GetUserEmojiCollectVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserEmojiCollectVersionsLogic) GetUserEmojiCollectVersions(in *emoji_rpc.GetUserEmojiCollectVersionsReq) (*emoji_rpc.GetUserEmojiCollectVersionsRes, error) {
	// todo: add your logic here and delete this line

	return &emoji_rpc.GetUserEmojiCollectVersionsRes{}, nil
}
