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

type GetSyncGroupMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的群成员版本
func NewGetSyncGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncGroupMembersLogic {
	return &GetSyncGroupMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncGroupMembersLogic) GetSyncGroupMembers(req *types.GetSyncGroupMembersReq) (resp *types.GetSyncGroupMembersRes, err error) {
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
		return &types.GetSyncGroupMembersRes{
			GroupVersions:   []types.GroupMembersVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	// 2. 获取变更的群成员
	serverTimestamp := time.Now().UnixMilli()

	memberResp, err := l.svcCtx.GroupRpc.GetGroupMembersListByIds(l.ctx, &group_rpc.GetGroupMembersListByIdsReq{
		GroupIDs: groupIDs,
		Since:    req.Since,
	})
	if err != nil {
		l.Errorf("获取变更的群成员失败: %v", err)
		return nil, err
	}

	// 3. 合并版本信息，按群组聚合最新版本
	groupVersionsMap := make(map[string]int64)
	for _, member := range memberResp.Members {
		if currentVersion, exists := groupVersionsMap[member.GroupID]; !exists || member.Version > currentVersion {
			groupVersionsMap[member.GroupID] = member.Version
		}
	}

	// 4. 转换为响应格式
	var groupVersions []types.GroupMembersVersionItem
	for groupID, version := range groupVersionsMap {
		groupVersions = append(groupVersions, types.GroupMembersVersionItem{
			GroupID: groupID,
			Version: version,
		})
	}

	return &types.GetSyncGroupMembersRes{
		GroupVersions:   groupVersions,
		ServerTimestamp: serverTimestamp,
	}, nil
}
