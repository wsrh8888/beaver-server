package logic

import (
	"context"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageContentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiPackageContentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageContentsLogic {
	return &GetEmojiPackageContentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiPackageContentsLogic) GetEmojiPackageContents(in *emoji_rpc.GetEmojiPackageContentsReq) (*emoji_rpc.GetEmojiPackageContentsRes, error) {
	// todo: add your logic here and delete this line

	return &emoji_rpc.GetEmojiPackageContentsRes{}, nil
}
