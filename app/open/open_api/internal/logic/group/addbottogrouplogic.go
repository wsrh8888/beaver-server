package group

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddBotToGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 将 Bot 加入群（Bot 进群后可接收群消息和 @提及）
func NewAddBotToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBotToGroupLogic {
	return &AddBotToGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddBotToGroupLogic) AddBotToGroup(req *types.AddBotToGroupReq) (resp *types.AddBotToGroupRes, err error) {
	// todo: add your logic here and delete this line

	return
}
