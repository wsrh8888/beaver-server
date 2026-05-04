// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package message

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendFileMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送文件消息
func NewSendFileMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendFileMessageLogic {
	return &SendFileMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendFileMessageLogic) SendFileMessage(req *types.SendFileMessageReq) (resp *types.SendFileMessageRes, err error) {
	// todo: add your logic here and delete this line

	return
}
