package logic

import (
	"context"
	"time"

	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户数据同步
func NewUserSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSyncLogic {
	return &UserSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserSyncLogic) UserSync(req *types.UserSyncReq) (resp *types.UserSyncRes, err error) {
	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	// 获取当前用户相关的用户ID列表（自己、好友、群友）
	relatedUserIDs, err := l.getRelatedUserIDs(req.UserID)
	if err != nil {
		l.Errorf("获取相关用户ID失败: %v", err)
		return nil, err
	}

	if len(relatedUserIDs) == 0 {
		l.Infof("用户没有相关用户数据需要同步，用户ID: %s", req.UserID)
		return &types.UserSyncRes{
			Users:       []types.UserSyncItem{},
			HasMore:     false,
			NextVersion: req.ToVersion + 1,
		}, nil
	}

	// 查询相关用户的版本号变化
	var users []user_models.UserModel
	err = l.svcCtx.DB.Where("uuid IN ? AND version > ? AND version <= ?",
		relatedUserIDs, req.FromVersion, req.ToVersion).
		Order("version ASC").
		Limit(limit + 1).
		Find(&users).Error
	if err != nil {
		l.Errorf("查询相关用户数据失败: %v", err)
		return nil, err
	}

	// 判断是否还有更多数据
	hasMore := len(users) > limit
	if hasMore {
		users = users[:limit]
	}

	// 转换为响应格式
	var userItems []types.UserSyncItem
	var nextVersion int64 = req.FromVersion

	for _, user := range users {
		userItems = append(userItems, types.UserSyncItem{
			UserID:   user.UUID,
			Nickname: user.NickName,
			Avatar:   user.Avatar,
			Abstract: user.Abstract,
			Phone:    user.Phone,
			Email:    user.Email,
			Gender:   user.Gender,
			Status:   user.Status,
			Version:  user.Version,
			CreateAt: time.Time(user.CreatedAt).Unix(),
			UpdateAt: time.Time(user.UpdatedAt).Unix(),
		})

		nextVersion = user.Version
	}

	if !hasMore {
		nextVersion = req.ToVersion + 1
	}

	resp = &types.UserSyncRes{
		Users:       userItems,
		HasMore:     hasMore,
		NextVersion: nextVersion,
	}

	l.Infof("用户数据同步完成，用户ID: %s, 相关用户数: %d, 返回用户数: %d, 还有更多: %v",
		req.UserID, len(relatedUserIDs), len(userItems), hasMore)
	return resp, nil
}

// getRelatedUserIDs 获取用户相关的用户ID列表（自己、好友、群友）
func (l *UserSyncLogic) getRelatedUserIDs(userID string) ([]string, error) {
	var relatedUserIDs []string

	// 1. 添加自己
	relatedUserIDs = append(relatedUserIDs, userID)

	// 2. 获取好友列表
	friendIDs, err := l.getFriendIDs(userID)
	if err != nil {
		l.Errorf("获取好友列表失败: %v", err)
		// 好友获取失败不影响同步，继续执行
	} else {
		relatedUserIDs = append(relatedUserIDs, friendIDs...)
	}

	// 3. 获取群友列表
	groupMemberIDs, err := l.getGroupMemberIDs(userID)
	if err != nil {
		l.Errorf("获取群成员列表失败: %v", err)
		// 群成员获取失败不影响同步，继续执行
	} else {
		relatedUserIDs = append(relatedUserIDs, groupMemberIDs...)
	}

	// 去重
	uniqueUserIDs := l.removeDuplicates(relatedUserIDs)

	l.Infof("获取相关用户ID列表，用户ID: %s, 总相关用户数: %d", userID, len(uniqueUserIDs))
	return uniqueUserIDs, nil
}

// getFriendIDs 获取用户的好友ID列表
func (l *UserSyncLogic) getFriendIDs(userID string) ([]string, error) {
	var friendIDs []string

	// 调用friend RPC获取好友列表
	response, err := l.svcCtx.FriendRpc.GetFriendIds(l.ctx, &friend_rpc.GetFriendIdsRequest{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	friendIDs = append(friendIDs, response.FriendIds...)

	l.Infof("获取好友列表成功，用户ID: %s, 好友数: %d", userID, len(friendIDs))
	return friendIDs, nil
}

// getGroupMemberIDs 获取用户群组中的其他成员ID列表
func (l *UserSyncLogic) getGroupMemberIDs(userID string) ([]string, error) {
	// 直接调用group RPC获取用户的所有群成员ID
	response, err := l.svcCtx.GroupRpc.GetUserGroupMembers(l.ctx, &group_rpc.GetUserGroupMembersReq{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	l.Infof("获取群成员列表成功，用户ID: %s, 群成员数: %d", userID, len(response.MemberIDs))
	return response.MemberIDs, nil
}

// removeDuplicates 去重
func (l *UserSyncLogic) removeDuplicates(userIDs []string) []string {
	seen := make(map[string]bool)
	var uniqueUserIDs []string

	for _, userID := range userIDs {
		if !seen[userID] {
			seen[userID] = true
			uniqueUserIDs = append(uniqueUserIDs, userID)
		}
	}

	return uniqueUserIDs
}
