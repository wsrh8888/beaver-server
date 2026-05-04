package event

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteEventSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除事件订阅
func NewDeleteEventSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEventSubscriptionLogic {
	return &DeleteEventSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteEventSubscriptionLogic) DeleteEventSubscription(req *types.DeleteEventSubscriptionReq) (resp *types.DeleteEventSubscriptionRes, err error) {
	// TODO: 删除事件订阅功能
	logx.Infof("删除事件订阅: subscriptionID=%d", req.SubscriptionID)

	return &types.DeleteEventSubscriptionRes{
		Success: true,
	}, nil
}
