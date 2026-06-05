package open

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOpenWebhookLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOpenWebhookLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOpenWebhookLogListLogic {
	return &GetOpenWebhookLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOpenWebhookLogListLogic) GetOpenWebhookLogList(req *types.GetOpenWebhookLogListReq) (resp *types.GetOpenWebhookLogListRes, err error) {
	rpcRes, err := l.svcCtx.OpenRpc.ListWebhookLogs(l.ctx, &open_rpc.ListWebhookLogsReq{
		Page:      int32(req.Page),
		PageSize:  int32(req.PageSize),
		AppId:     req.AppID,
		EventType: req.EventType,
		Status:    int32(req.Status),
	})
	if err != nil {
		l.Errorf("查询 Webhook 日志失败: %v", err)
		return nil, err
	}

	list := make([]types.OpenWebhookLogInfo, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		list = append(list, types.OpenWebhookLogInfo{
			ID:           item.Id,
			AppID:        item.AppId,
			EventID:      item.EventId,
			EventType:    item.EventType,
			TargetURL:    item.TargetUrl,
			HTTPStatus:   int(item.HttpStatus),
			LatencyMs:    item.LatencyMs,
			RetryCount:   int(item.RetryCount),
			Status:       int(item.Status),
			ErrorMessage: item.ErrorMessage,
			CreatedAt:    item.CreatedAt,
		})
	}

	return &types.GetOpenWebhookLogListRes{Total: rpcRes.Total, List: list}, nil
}
