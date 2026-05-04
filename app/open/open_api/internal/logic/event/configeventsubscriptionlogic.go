package event

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigEventSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 配置事件订阅
func NewConfigEventSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigEventSubscriptionLogic {
	return &ConfigEventSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigEventSubscriptionLogic) ConfigEventSubscription(req *types.ConfigEventSubscriptionReq) (resp *types.ConfigEventSubscriptionRes, err error) {
	// TODO: 事件订阅配置功能
	logx.Infof("配置事件订阅: appID=%s, eventType=%s", req.AppID, req.EventType)

	return &types.ConfigEventSubscriptionRes{
		SubscriptionID: 1, // 简化处理，实际应该返回创建的订阅 ID
	}, nil
}
