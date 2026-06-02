package event

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEventSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新事件订阅
func NewUpdateEventSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEventSubscriptionLogic {
	return &UpdateEventSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEventSubscriptionLogic) UpdateEventSubscription(req *types.UpdateEventSubscriptionReq) (resp *types.UpdateEventSubscriptionRes, err error) {
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// todo: add your logic here and delete this line

	return
}
