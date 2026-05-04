package message

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCardMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新卡片消息
func NewUpdateCardMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCardMessageLogic {
	return &UpdateCardMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCardMessageLogic) UpdateCardMessage(req *types.UpdateCardMessageReq) (resp *types.UpdateCardMessageRes, err error) {
	// TODO: 需要调用 Chat RPC 更新卡片消息
	logx.Infof("更新卡片消息: messageID=%s", req.MessageID)

	return &types.UpdateCardMessageRes{
		Success: true,
	}, nil
}
