package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/notification/notification_rpc/types/notification_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncNotificationInboxesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取通知收件箱版本摘要
func NewGetSyncNotificationInboxesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncNotificationInboxesLogic {
	return &GetSyncNotificationInboxesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncNotificationInboxesLogic) GetSyncNotificationInboxes(req *types.GetSyncNotificationInboxesReq) (resp *types.GetSyncNotificationInboxesRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	rpcResp, err := l.svcCtx.NotificationRpc.GetInboxVersions(l.ctx, &notification_rpc.GetInboxVersionsReq{
		UserId:       req.UserID,
		SinceVersion: req.SinceVersion,
		Limit:        req.Limit,
	})
	if err != nil {
		l.Errorf("获取通知收件箱版本摘要失败: userId=%s, sinceVersion=%d, limit=%d, err=%v", req.UserID, req.SinceVersion, req.Limit, err)
		return nil, err
	}

	inboxVersions := make([]types.NotificationInboxVersionItem, 0)
	for _, item := range rpcResp.InboxVersions {
		inboxVersions = append(inboxVersions, types.NotificationInboxVersionItem{
			EventID: item.EventId,
			Version: item.Version,
		})
	}

	return &types.GetSyncNotificationInboxesRes{
		InboxVersions:  inboxVersions,
		MaxVersion:     rpcResp.MaxVersion,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}

