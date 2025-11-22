package logic

import (
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncGroupRequestsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的入群申请版本
func NewGetSyncGroupRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncGroupRequestsLogic {
	return &GetSyncGroupRequestsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncGroupRequestsLogic) GetSyncGroupRequests(req *types.GetSyncGroupRequestsReq) (resp *types.GetSyncGroupRequestsRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	// 获取用户群组申请版本信息
	versionResp, err := l.svcCtx.GroupRpc.GetUserGroupRequestVersions(l.ctx, &group_rpc.GetUserGroupRequestVersionsReq{
		UserID: userId,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取用户群组申请版本失败: %v", err)
		return nil, err
	}

	// 转换为响应格式，确保返回空数组而不是null
	groupVersions := make([]types.GroupRequestsVersionItem, 0)
	if versionResp.Versions != nil {
		for _, version := range versionResp.Versions {
			groupVersions = append(groupVersions, types.GroupRequestsVersionItem{
				GroupID: version.GroupID,
				Version: version.Version,
			})
		}
	}

	return &types.GetSyncGroupRequestsRes{
		GroupVersions:   groupVersions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
