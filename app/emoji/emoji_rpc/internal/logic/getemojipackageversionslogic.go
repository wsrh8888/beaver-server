package logic

import (
	"context"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiPackageVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageVersionsLogic {
	return &GetEmojiPackageVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiPackageVersionsLogic) GetEmojiPackageVersions(in *emoji_rpc.GetEmojiPackageVersionsReq) (*emoji_rpc.GetEmojiPackageVersionsRes, error) {
	// todo: add your logic here and delete this line

	return &emoji_rpc.GetEmojiPackageVersionsRes{}, nil
}
