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

type GetSyncAllGroupsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的群组版本信息
func NewGetSyncAllGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncAllGroupsLogic {
	return &GetSyncAllGroupsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncAllGroupsLogic) GetSyncAllGroups(req *types.GetSyncAllGroupsReq) (resp *types.GetSyncAllGroupsRes, err error) {
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
		return &types.GetSyncAllGroupsRes{
			GroupVersions:   []types.GroupVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	// 2. 获取变更的群组资料、群成员、入群申请
	serverTimestamp := time.Now().UnixMilli()

	// 获取变更的群组资料
	groupResp, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{
		GroupIDs: groupIDs,
		Since:    req.Since,
	})
	if err != nil {
		l.Errorf("获取变更的群组资料失败: %v", err)
		return nil, err
	}

	// 获取变更的群成员
	memberResp, err := l.svcCtx.GroupRpc.GetGroupMembersListByIds(l.ctx, &group_rpc.GetGroupMembersListByIdsReq{
		GroupIDs: groupIDs,
		Since:    req.Since,
	})
	if err != nil {
		l.Errorf("获取变更的群成员失败: %v", err)
		return nil, err
	}

	// 获取变更的入群申请
	requestResp, err := l.svcCtx.GroupRpc.GetGroupJoinRequestsListByIds(l.ctx, &group_rpc.GetGroupJoinRequestsListByIdsReq{
		GroupIDs: groupIDs,
		Since:    req.Since,
	})
	if err != nil {
		l.Errorf("获取变更的入群申请失败: %v", err)
		return nil, err
	}

	// 3. 合并版本信息
	groupVersionsMap := make(map[string]*types.GroupVersionItem)

	// 处理群组资料版本
	for _, group := range groupResp.Groups {
		groupVersionsMap[group.GroupID] = &types.GroupVersionItem{
			GroupID:        group.GroupID,
			GroupVersion:   group.Version,
			MemberVersion:  0, // 稍后更新
			RequestVersion: 0, // 稍后更新
		}
	}

	// 处理群成员版本
	for _, member := range memberResp.Members {
		if item, exists := groupVersionsMap[member.GroupID]; exists {
			if member.Version > item.MemberVersion {
				item.MemberVersion = member.Version
			}
		} else {
			groupVersionsMap[member.GroupID] = &types.GroupVersionItem{
				GroupID:        member.GroupID,
				GroupVersion:   0,
				MemberVersion:  member.Version,
				RequestVersion: 0,
			}
		}
	}

	// 处理入群申请版本
	for _, request := range requestResp.Requests {
		if item, exists := groupVersionsMap[request.GroupID]; exists {
			if request.Version > item.RequestVersion {
				item.RequestVersion = request.Version
			}
		} else {
			groupVersionsMap[request.GroupID] = &types.GroupVersionItem{
				GroupID:        request.GroupID,
				GroupVersion:   0,
				MemberVersion:  0,
				RequestVersion: request.Version,
			}
		}
	}

	// 转换为响应格式
	var groupVersions []types.GroupVersionItem
	for _, item := range groupVersionsMap {
		groupVersions = append(groupVersions, *item)
	}

	return &types.GetSyncAllGroupsRes{
		GroupVersions:   groupVersions,
		ServerTimestamp: serverTimestamp,
	}, nil
}
