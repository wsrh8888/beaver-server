package logic

import (
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/moment/moment_rpc/types/moment_rpc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncMomentCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的动态评论版本
func NewGetSyncMomentCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncMomentCommentsLogic {
	return &GetSyncMomentCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncMomentCommentsLogic) GetSyncMomentComments(req *types.GetSyncMomentCommentsReq) (resp *types.GetSyncMomentCommentsRes, err error) {
	// 调用moment RPC服务获取版本摘要

	rpcReq := &moment_rpc.GetMomentCommentVersionsReq{
		UserId: req.UserID,
		Since:  req.Since,
	}

	rpcResp, err := l.svcCtx.MomentRpc.GetMomentCommentVersions(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("调用Moment服务的GetMomentCommentVersions失败: %v", err)
		return &types.GetSyncMomentCommentsRes{
			MomentCommentVersions: []types.MomentCommentVersionItem{},
			ServerTimestamp:       0,
		}, nil
	}

	// 转换响应格式，确保返回空数组而不是null
	versions := make([]types.MomentCommentVersionItem, 0)
	if rpcResp.MomentCommentVersions != nil {
		for _, v := range rpcResp.MomentCommentVersions {
			versions = append(versions, types.MomentCommentVersionItem{
				UserID:  v.UserId,
				Version: v.Version,
			})
		}
	}

	return &types.GetSyncMomentCommentsRes{
		MomentCommentVersions: versions,
		ServerTimestamp:       rpcResp.ServerTimestamp,
	}, nil
}
