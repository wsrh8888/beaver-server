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

type GetSyncNotificationReadCursorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取通知已读游标版本摘要
func NewGetSyncNotificationReadCursorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncNotificationReadCursorsLogic {
	return &GetSyncNotificationReadCursorsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncNotificationReadCursorsLogic) GetSyncNotificationReadCursors(req *types.GetSyncNotificationReadCursorsReq) (resp *types.GetSyncNotificationReadCursorsRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	rpcResp, err := l.svcCtx.NotificationRpc.GetReadCursorVersions(l.ctx, &notification_rpc.GetReadCursorVersionsReq{
		UserId:       req.UserID,
		SinceVersion: req.SinceVersion,
	})
	if err != nil {
		l.Errorf("获取通知已读游标版本摘要失败: userId=%s, sinceVersion=%d, err=%v", req.UserID, req.SinceVersion, err)
		return nil, err
	}

	cursorVersions := make([]types.NotificationReadCursorVersionItem, 0)
	for _, item := range rpcResp.CursorVersions {
		cursorVersions = append(cursorVersions, types.NotificationReadCursorVersionItem{
			Category: item.Category,
			Version:  item.Version,
		})
	}

	return &types.GetSyncNotificationReadCursorsRes{
		CursorVersions: cursorVersions,
		MaxVersion:     rpcResp.MaxVersion,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}

