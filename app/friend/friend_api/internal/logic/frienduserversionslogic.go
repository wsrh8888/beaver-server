package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/list_query"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendUserVersionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取好友用户版本信息（用于数据同步）
func NewFriendUserVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendUserVersionsLogic {
	return &FriendUserVersionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendUserVersionsLogic) FriendUserVersions(req *types.FriendUserVersionsReq) (resp *types.FriendUserVersionsRes, err error) {
	// 参数验证
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	// 设置默认分页参数
	if req.Limit <= 0 {
		req.Limit = 100
	}

	// 查询所有好友关系（不分页，因为需要获取所有好友的版本信息）
	allFriends, _, _ := list_query.ListQuery(l.svcCtx.DB, friend_models.FriendModel{}, list_query.Option{
		Where: l.svcCtx.DB.Where("(send_user_id = ? OR rev_user_id = ?) AND is_deleted = ?", req.UserID, req.UserID, false),
	})

	// 收集所有用户ID（好友 + 当前用户）
	userIdSet := make(map[string]bool)
	var allUserIds []string

	// 始终添加当前用户（version为0，强制同步）
	allUserIds = append(allUserIds, req.UserID)
	userIdSet[req.UserID] = true

	for _, friend := range allFriends {
		if friend.SendUserID != req.UserID && friend.SendUserID != "" && !userIdSet[friend.SendUserID] {
			allUserIds = append(allUserIds, friend.SendUserID)
			userIdSet[friend.SendUserID] = true
		}
		if friend.RevUserID != req.UserID && friend.RevUserID != "" && !userIdSet[friend.RevUserID] {
			allUserIds = append(allUserIds, friend.RevUserID)
			userIdSet[friend.RevUserID] = true
		}
	}

	// 应用分页
	totalUsers := len(allUserIds)
	start := req.Offset
	end := req.Offset + req.Limit
	if start >= totalUsers {
		// 超出范围，返回空结果
		return &types.FriendUserVersionsRes{
			UserVersions: []types.FriendUserVersionItem{},
			Total:        totalUsers,
		}, nil
	}
	if end > totalUsers {
		end = totalUsers
	}

	pagedUserIds := allUserIds[start:end]

	// 如果没有用户，直接返回
	if len(pagedUserIds) == 0 {
		return &types.FriendUserVersionsRes{
			UserVersions: []types.FriendUserVersionItem{},
			Total:        totalUsers,
		}, nil
	}

	// 调用用户RPC服务获取用户版本信息
	userVersionsResp, err := l.svcCtx.UserRpc.UserVersions(l.ctx, &user_rpc.UserVersionsReq{
		UserIds: pagedUserIds,
	})
	if err != nil {
		l.Logger.Errorf("获取用户版本信息失败: %v", err)
		return nil, errors.New("获取用户版本信息失败")
	}

	// 构造响应
	var userVersions []types.FriendUserVersionItem
	for _, userId := range pagedUserIds {
		version := int64(1) // 默认版本号
		if userVersion, exists := userVersionsResp.UserVersions[userId]; exists {
			version = userVersion
		}

		// 当前用户始终返回version=0，强制同步
		if userId == req.UserID {
			version = 0
		}

		userVersions = append(userVersions, types.FriendUserVersionItem{
			UserID:  userId,
			Version: version,
		})
	}

	l.Logger.Infof("获取用户版本信息成功: userID=%s, total=%d, returned=%d", req.UserID, totalUsers, len(userVersions))

	return &types.FriendUserVersionsRes{
		UserVersions: userVersions,
		Total:        totalUsers,
	}, nil
}
