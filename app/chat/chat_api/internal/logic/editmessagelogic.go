package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditMessageLogic {
	return &EditMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditMessageLogic) EditMessage(req *types.EditMessageReq) (resp *types.EditMessageRes, err error) {
	// todo: add your logic here and delete this line

	return
}
