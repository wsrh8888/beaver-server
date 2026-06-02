package event

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEventSubscriptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取事件订阅列表
func NewListEventSubscriptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEventSubscriptionsLogic {
	return &ListEventSubscriptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListEventSubscriptionsLogic) ListEventSubscriptions(req *types.ListEventSubscriptionsReq) (resp *types.ListEventSubscriptionsRes, err error) {
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// todo: add your logic here and delete this line

	return
}
