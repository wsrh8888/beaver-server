package logic

import (
	"context"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserEmojiPackageCollectVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserEmojiPackageCollectVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserEmojiPackageCollectVersionsLogic {
	return &GetUserEmojiPackageCollectVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserEmojiPackageCollectVersionsLogic) GetUserEmojiPackageCollectVersions(in *emoji_rpc.GetUserEmojiPackageCollectVersionsReq) (*emoji_rpc.GetUserEmojiPackageCollectVersionsRes, error) {
	// todo: add your logic here and delete this line

	return &emoji_rpc.GetUserEmojiPackageCollectVersionsRes{}, nil
}
