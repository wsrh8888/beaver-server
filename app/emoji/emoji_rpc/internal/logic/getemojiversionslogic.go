package logic

import (
	"context"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiVersionsLogic {
	return &GetEmojiVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 版本摘要查询接口（为微服务架构设计）
func (l *GetEmojiVersionsLogic) GetEmojiVersions(in *emoji_rpc.GetEmojiVersionsReq) (*emoji_rpc.GetEmojiVersionsRes, error) {
	// todo: add your logic here and delete this line

	return &emoji_rpc.GetEmojiVersionsRes{}, nil
}
