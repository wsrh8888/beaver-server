package logic

import (
	"context"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojisByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojisByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojisByIdsLogic {
	return &GetEmojisByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 基础数据查询接口
func (l *GetEmojisByIdsLogic) GetEmojisByIds(in *emoji_rpc.GetEmojisByIdsReq) (*emoji_rpc.GetEmojisByIdsRes, error) {
	// todo: add your logic here and delete this line

	return &emoji_rpc.GetEmojisByIdsRes{}, nil
}
