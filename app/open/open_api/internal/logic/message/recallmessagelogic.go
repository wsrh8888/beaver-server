package message

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecallMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 撤回消息
func NewRecallMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecallMessageLogic {
	return &RecallMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecallMessageLogic) RecallMessage(req *types.RecallMessageReq) (resp *types.RecallMessageRes, err error) {
	// TODO: 需要调用 Chat RPC 撤回消息
	logx.Infof("撤回消息: messageID=%s", req.MessageID)

	return &types.RecallMessageRes{
		Success: true,
	}, nil
}
