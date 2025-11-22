package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncAllUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取需要同步的用户列表
func NewGetSyncAllUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncAllUsersLogic {
	return &GetSyncAllUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncAllUsersLogic) GetSyncAllUsers(req *types.GetSyncAllUsersReq) (resp *types.GetSyncAllUsersRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	// 根据类型获取相关用户ID列表
	var relatedUserIds []string

	// 确定需要获取哪些数据
	needFriends := req.Type == "" || req.Type == "all" || req.Type == "friends"
	needGroups := req.Type == "" || req.Type == "all" || req.Type == "groups"

	// 使用map进行去重
	userMap := make(map[string]bool)

	// 始终包含自己的ID
	userMap[userId] = true

	// 获取好友列表
	if needFriends {
		friendResp, err := l.svcCtx.FriendRpc.GetFriendIds(l.ctx, &friend_rpc.GetFriendIdsRequest{
			UserID: userId,
		})
		if err != nil {
			l.Errorf("获取好友列表失败: %v", err)
			return nil, err
		}
		for _, uid := range friendResp.FriendIds {
			userMap[uid] = true
		}
	}

	// 获取群成员列表
	if needGroups {
		groupResp, err := l.svcCtx.GroupRpc.GetUserGroupMembers(l.ctx, &group_rpc.GetUserGroupMembersReq{
			UserID: userId,
		})
		if err != nil {
			l.Errorf("获取群成员列表失败: %v", err)
			return nil, err
		}
		for _, uid := range groupResp.MemberIDs {
			userMap[uid] = true
		}
	}

	// 获取好友验证申请相关的用户（申请加我好友的人）
	if needFriends {
		verifyResp, err := l.svcCtx.FriendRpc.GetFriendVerifyVersions(l.ctx, &friend_rpc.GetFriendVerifyVersionsReq{
			UserId: userId,
			Since:  0, // 获取所有相关验证，不需要增量
		})
		if err != nil {
			l.Errorf("获取好友验证列表失败: %v", err)
			return nil, err
		}

		// 获取验证详情以提取申请者ID
		if len(verifyResp.FriendVerifyVersions) > 0 {
			uuids := make([]string, 0, len(verifyResp.FriendVerifyVersions))
			for _, verify := range verifyResp.FriendVerifyVersions {
				uuids = append(uuids, verify.Uuid)
			}

			verifyDetailResp, err := l.svcCtx.FriendRpc.GetFriendVerifiesListByIds(l.ctx, &friend_rpc.GetFriendVerifiesListByIdsReq{
				Uuids: uuids,
			})
			if err != nil {
				l.Errorf("获取好友验证详情失败: %v", err)
				return nil, err
			}

			for _, verify := range verifyDetailResp.FriendVerifies {
				// 添加发送申请的用户
				userMap[verify.SendUserId] = true
				// 添加接收申请的用户（通常是当前用户自己，但为了完整性也添加）
				userMap[verify.RevUserId] = true
			}
		}
	}

	// 获取群申请相关的用户（申请加入我管理的群的人）
	if needGroups {
		// 1. 获取用户管理的群组
		groupIDsResp, err := l.svcCtx.GroupRpc.GetUserGroupIDs(l.ctx, &group_rpc.GetUserGroupIDsReq{
			UserID: userId,
		})
		if err != nil {
			l.Errorf("获取用户群组ID列表失败: %v", err)
			return nil, err
		}

		if len(groupIDsResp.GroupIDs) > 0 {
			// 2. 获取这些群的加入申请
			requestResp, err := l.svcCtx.GroupRpc.GetGroupJoinRequestsListByIds(l.ctx, &group_rpc.GetGroupJoinRequestsListByIdsReq{
				GroupIDs: groupIDsResp.GroupIDs,
				Since:    0, // 获取所有申请，不需要增量
			})
			if err != nil {
				l.Errorf("获取群加入申请失败: %v", err)
				return nil, err
			}

			// 3. 提取申请者用户ID
			for _, request := range requestResp.Requests {
				userMap[request.UserID] = true
			}
		}
	}

	// 检查是否有不支持的类型
	if !needFriends && !needGroups {
		l.Errorf("不支持的类型: %s", req.Type)
		return nil, errors.New("不支持的类型")
	}

	// 转换为切片
	for uid := range userMap {
		relatedUserIds = append(relatedUserIds, uid)
	}

	if len(relatedUserIds) == 0 {
		return &types.GetSyncAllUsersRes{
			UserVersions:    []types.UserVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	// 调用user RPC获取用户信息，支持增量同步
	userResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
		UserIdList:     relatedUserIds,
		SinceTimestamp: req.Since, // 传递since参数进行增量过滤
	})
	if err != nil {
		l.Errorf("调用user RPC获取用户信息失败: %v", err)
		return nil, err
	}

	// 转换为版本摘要，确保返回空数组而不是null
	userVersions := make([]types.UserVersionItem, 0)
	if userResp.UserInfo != nil {
		for userId, userInfo := range userResp.UserInfo {
			userVersions = append(userVersions, types.UserVersionItem{
				UserID:  userId,
				Version: userInfo.Version,
			})
		}
	}

	return &types.GetSyncAllUsersRes{
		UserVersions:    userVersions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
