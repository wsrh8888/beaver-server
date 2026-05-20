package event

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWebhookListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前应用的 Webhook 订阅列表
func NewGetWebhookListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWebhookListLogic {
	return &GetWebhookListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWebhookListLogic) GetWebhookList(req *types.GetWebhookListReq) (resp *types.GetWebhookListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
