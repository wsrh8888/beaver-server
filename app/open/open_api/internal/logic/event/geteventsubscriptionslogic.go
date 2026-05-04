package event

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventSubscriptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取事件订阅列表
func NewGetEventSubscriptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventSubscriptionsLogic {
	return &GetEventSubscriptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventSubscriptionsLogic) GetEventSubscriptions(req *types.GetEventSubscriptionsReq) (resp *types.GetEventSubscriptionsRes, err error) {
	// TODO: 获取事件订阅列表
	logx.Infof("获取事件订阅列表: appID=%s", req.AppID)

	return &types.GetEventSubscriptionsRes{
		Total: 0,
		List:  []types.EventSubscriptionInfo{},
	}, nil
}
