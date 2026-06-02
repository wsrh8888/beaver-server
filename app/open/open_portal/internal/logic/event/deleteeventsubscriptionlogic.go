package event

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

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
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// todo: add your logic here and delete this line

	return
}
