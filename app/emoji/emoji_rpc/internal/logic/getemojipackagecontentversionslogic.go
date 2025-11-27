package logic

import (
	"context"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageContentVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiPackageContentVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageContentVersionsLogic {
	return &GetEmojiPackageContentVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiPackageContentVersionsLogic) GetEmojiPackageContentVersions(in *emoji_rpc.GetEmojiPackageContentVersionsReq) (*emoji_rpc.GetEmojiPackageContentVersionsRes, error) {
	// todo: add your logic here and delete this line

	return &emoji_rpc.GetEmojiPackageContentVersionsRes{}, nil
}
