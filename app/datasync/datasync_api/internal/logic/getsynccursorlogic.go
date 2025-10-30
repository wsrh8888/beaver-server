package logic

import (
	"context"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/datasync/datasync_rpc/types/types/datasync_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncCursorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取同步游标
func NewGetSyncCursorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncCursorLogic {
	return &GetSyncCursorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncCursorLogic) GetSyncCursor(req *types.GetSyncCursorReq) (resp *types.GetSyncCursorRes, err error) {
	// 调用 RPC 服务获取同步游标
	rpcResp, err := l.svcCtx.DatasyncRpc.GetSyncCursor(l.ctx, &datasync_rpc.GetSyncCursorReq{
		UserId:   req.UserID,
		DeviceId: req.DeviceID,
		DataType: req.DataType,
	})

	if err != nil {
		l.Errorf("调用 RPC 获取同步游标失败: %v", err)
		return nil, err
	}

	return &types.GetSyncCursorRes{
		DataType: req.DataType,
		LastSeq:  rpcResp.ServerLatest, // 使用ServerLatest作为服务端最新序列号
	}, nil
}
