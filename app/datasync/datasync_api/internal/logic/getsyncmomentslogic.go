package logic

import (
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/moment/moment_rpc/types/moment_rpc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncMomentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的动态版本
func NewGetSyncMomentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncMomentsLogic {
	return &GetSyncMomentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncMomentsLogic) GetSyncMoments(req *types.GetSyncMomentsReq) (resp *types.GetSyncMomentsRes, err error) {
	// 调用moment RPC服务获取版本摘要
	rpcReq := &moment_rpc.GetMomentVersionsReq{
		UserId: req.UserID,
		Since:  req.Since,
	}

	rpcResp, err := l.svcCtx.MomentRpc.GetMomentVersions(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("调用Moment RPC服务失败: %v", err)
		return &types.GetSyncMomentsRes{
			MomentVersions:  []types.MomentVersionItem{},
			ServerTimestamp: 0,
		}, nil
	}

	// 转换响应格式
	var versions []types.MomentVersionItem
	for _, v := range rpcResp.MomentVersions {
		versions = append(versions, types.MomentVersionItem{
			UserID:  v.UserId,
			Version: v.Version,
		})
	}

	return &types.GetSyncMomentsRes{
		MomentVersions:  versions,
		ServerTimestamp: rpcResp.ServerTimestamp,
	}, nil
}
