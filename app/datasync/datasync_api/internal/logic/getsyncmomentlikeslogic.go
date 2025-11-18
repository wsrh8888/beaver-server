package logic

import (
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/moment/moment_rpc/types/moment_rpc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncMomentLikesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的动态点赞版本
func NewGetSyncMomentLikesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncMomentLikesLogic {
	return &GetSyncMomentLikesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncMomentLikesLogic) GetSyncMomentLikes(req *types.GetSyncMomentLikesReq) (resp *types.GetSyncMomentLikesRes, err error) {
	// 调用moment RPC服务获取版本摘要
	rpcReq := &moment_rpc.GetMomentLikeVersionsReq{
		UserId: req.UserID,
		Since:  req.Since,
	}

	rpcResp, err := l.svcCtx.MomentRpc.GetMomentLikeVersions(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("调用Moment RPC服务失败: %v", err)
		return &types.GetSyncMomentLikesRes{
			MomentLikeVersions: []types.MomentLikeVersionItem{},
			ServerTimestamp:    0,
		}, nil
	}

	// 转换响应格式
	var versions []types.MomentLikeVersionItem
	for _, v := range rpcResp.MomentLikeVersions {
		versions = append(versions, types.MomentLikeVersionItem{
			UserID:  v.UserId,
			Version: v.Version,
		})
	}

	return &types.GetSyncMomentLikesRes{
		MomentLikeVersions: versions,
		ServerTimestamp:    rpcResp.ServerTimestamp,
	}, nil
}
