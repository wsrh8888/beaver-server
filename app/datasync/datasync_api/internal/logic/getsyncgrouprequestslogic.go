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

	// 1. 获取用户加入的群组ID列表
	groupIDsResp, err := l.svcCtx.GroupRpc.GetUserGroupIDs(l.ctx, &group_rpc.GetUserGroupIDsReq{
		UserID: userId,
	})
	if err != nil {
		l.Errorf("获取用户群组ID列表失败: %v", err)
		return nil, err
	}

	groupIDs := groupIDsResp.GroupIDs
	if len(groupIDs) == 0 {
		return &types.GetSyncGroupRequestsRes{
			GroupVersions:   []types.GroupRequestsVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	// 2. 获取变更的入群申请
	serverTimestamp := time.Now().UnixMilli()

	requestResp, err := l.svcCtx.GroupRpc.GetGroupJoinRequestsListByIds(l.ctx, &group_rpc.GetGroupJoinRequestsListByIdsReq{
		GroupIDs: groupIDs,
		Since:    req.Since,
	})
	if err != nil {
		l.Errorf("获取变更的入群申请失败: %v", err)
		return nil, err
	}

	// 3. 合并版本信息，按群组聚合最新版本
	groupVersionsMap := make(map[string]int64)
	for _, request := range requestResp.Requests {
		if currentVersion, exists := groupVersionsMap[request.GroupID]; !exists || request.Version > currentVersion {
			groupVersionsMap[request.GroupID] = request.Version
		}
	}

	// 4. 转换为响应格式
	var groupVersions []types.GroupRequestsVersionItem
	for groupID, version := range groupVersionsMap {
		groupVersions = append(groupVersions, types.GroupRequestsVersionItem{
			GroupID: groupID,
			Version: version,
		})
	}

	return &types.GetSyncGroupRequestsRes{
		GroupVersions:   groupVersions,
		ServerTimestamp: serverTimestamp,
	}, nil
}
