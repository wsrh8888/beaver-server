package logic

import (
	"context"
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/notification/notification_rpc/types/notification_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncNotificationEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取通知事件版本摘要
func NewGetSyncNotificationEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncNotificationEventsLogic {
	return &GetSyncNotificationEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncNotificationEventsLogic) GetSyncNotificationEvents(req *types.GetSyncNotificationEventsReq) (resp *types.GetSyncNotificationEventsRes, err error) {
	rpcResp, err := l.svcCtx.NotificationRpc.GetEventVersions(l.ctx, &notification_rpc.GetEventVersionsReq{
		SinceVersion: req.SinceVersion,
		Limit:        req.Limit,
	})
	if err != nil {
		l.Errorf("获取通知事件版本摘要失败: sinceVersion=%d, limit=%d, err=%v", req.SinceVersion, req.Limit, err)
		return nil, err
	}

	eventVersions := make([]types.NotificationEventVersionItem, 0)
	for _, item := range rpcResp.EventVersions {
		eventVersions = append(eventVersions, types.NotificationEventVersionItem{
			EventID: item.EventId,
			Version: item.Version,
		})
	}

	return &types.GetSyncNotificationEventsRes{
		EventVersions:   eventVersions,
		MaxVersion:      rpcResp.MaxVersion,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}

