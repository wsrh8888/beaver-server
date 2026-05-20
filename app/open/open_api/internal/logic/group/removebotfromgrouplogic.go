package group

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveBotFromGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 将 Bot 从群移除
func NewRemoveBotFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveBotFromGroupLogic {
	return &RemoveBotFromGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveBotFromGroupLogic) RemoveBotFromGroup(req *types.RemoveBotFromGroupReq) (resp *types.RemoveBotFromGroupRes, err error) {
	// todo: add your logic here and delete this line

	return
}
