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
	// TODO: 需要调用 Chat RPC 发送文件消息
	logx.Infof("发送文件消息: targetID=%s, fileName=%s", req.TargetID, req.FileName)

	return &types.SendFileMessageRes{
		MessageID: "msg_file_xxx",
	}, nil
}
