package logic

import (
	"context"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackagesByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiPackagesByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackagesByIdsLogic {
	return &GetEmojiPackagesByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiPackagesByIdsLogic) GetEmojiPackagesByIds(in *emoji_rpc.GetEmojiPackagesByIdsReq) (*emoji_rpc.GetEmojiPackagesByIdsRes, error) {
	// todo: add your logic here and delete this line

	return &emoji_rpc.GetEmojiPackagesByIdsRes{}, nil
}
