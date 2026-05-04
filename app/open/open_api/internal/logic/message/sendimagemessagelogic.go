package message

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendImageMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送图片消息
func NewSendImageMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendImageMessageLogic {
	return &SendImageMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendImageMessageLogic) SendImageMessage(req *types.SendImageMessageReq) (resp *types.SendImageMessageRes, err error) {
	// TODO: 需要调用 Chat RPC 发送图片消息
	logx.Infof("发送图片消息: targetID=%s, imageURL=%s", req.TargetID, req.ImageURL)

	return &types.SendImageMessageRes{
		MessageID: "msg_image_xxx",
	}, nil
}
