package logic

import (
	"context"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/datasync/datasync_rpc/types/types/datasync_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSyncCursorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新同步游标
func NewUpdateSyncCursorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSyncCursorLogic {
	return &UpdateSyncCursorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSyncCursorLogic) UpdateSyncCursor(req *types.UpdateSyncCursorReq) (resp *types.UpdateSyncCursorRes, err error) {
	// 调用 RPC 服务更新同步游标
	_, err = l.svcCtx.DatasyncRpc.UpdateSyncCursor(l.ctx, &datasync_rpc.UpdateSyncCursorReq{
		UserId:         req.UserID,
		DeviceId:       req.DeviceID,
		DataType:       req.DataType,
		LastSeq:        req.LastSeq, // 使用LastSeq替代LastSyncTime
		ConversationId: req.ConversationID,
	})

	if err != nil {
		l.Errorf("调用 RPC 更新同步游标失败: %v", err)
		return nil, err
	}

	return &types.UpdateSyncCursorRes{}, nil
}
