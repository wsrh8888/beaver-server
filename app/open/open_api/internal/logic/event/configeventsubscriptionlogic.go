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
	// todo: add your logic here and delete this line

	return
}
