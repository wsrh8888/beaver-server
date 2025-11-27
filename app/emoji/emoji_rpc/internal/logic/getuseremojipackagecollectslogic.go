package logic

import (
	"context"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserEmojiPackageCollectsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserEmojiPackageCollectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserEmojiPackageCollectsLogic {
	return &GetUserEmojiPackageCollectsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserEmojiPackageCollectsLogic) GetUserEmojiPackageCollects(in *emoji_rpc.GetUserEmojiPackageCollectsReq) (*emoji_rpc.GetUserEmojiPackageCollectsRes, error) {
	// todo: add your logic here and delete this line

	return &emoji_rpc.GetUserEmojiPackageCollectsRes{}, nil
}
