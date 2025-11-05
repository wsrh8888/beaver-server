package logic

import (
	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserGroupVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserGroupVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserGroupVersionsLogic {
	return &GetUserGroupVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserGroupVersionsLogic) GetUserGroupVersions(req *types.GetUserGroupVersionsReq) (resp *types.GetUserGroupVersionsRes, err error) {
	// 获取用户加入的所有群组及对应版本信息
	var userGroupMembers []group_models.GroupMemberModel
	err = l.svcCtx.DB.Where("user_id = ? AND status = 1", req.UserID).
		Select("group_id").
		Find(&userGroupMembers).Error
	if err != nil {
		l.Errorf("查询用户群组成员关系失败: %v", err)
		return nil, err
	}

	if len(userGroupMembers) == 0 {
		return &types.GetUserGroupVersionsRes{
			Groups: []types.GroupVersionItem{},
		}, nil
	}

	// 提取群组ID列表
	groupIDs := make([]string, len(userGroupMembers))
	for i, member := range userGroupMembers {
		groupIDs[i] = member.GroupID
	}

	// 获取每个群组的最新版本信息
	var groups []group_models.GroupModel
	err = l.svcCtx.DB.Where("group_id IN (?)", groupIDs).Find(&groups).Error
	if err != nil {
		l.Errorf("查询群组信息失败: %v", err)
		return nil, err
	}

	// 为每个群组构建版本信息
	groupVersions := make(map[string]*types.GroupVersionItem)
	for _, group := range groups {
		// 获取群成员最新版本
		var memberVersion int64
		err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
			Where("group_id = ?", group.GroupID).
			Select("COALESCE(MAX(version), 0)").
			Scan(&memberVersion).Error
		if err != nil {
			l.Errorf("查询群成员版本失败, groupId: %s, error: %v", group.GroupID, err)
			memberVersion = 0
		}

		// 获取入群申请最新版本
		var requestVersion int64
		err = l.svcCtx.DB.Model(&group_models.GroupJoinRequestModel{}).
			Where("group_id = ?", group.GroupID).
			Select("COALESCE(MAX(version), 0)").
			Scan(&requestVersion).Error
		if err != nil {
			l.Errorf("查询入群申请版本失败, groupId: %s, error: %v", group.GroupID, err)
			requestVersion = 0
		}

		groupVersions[group.GroupID] = &types.GroupVersionItem{
			GroupID:        group.GroupID,
			GroupVersion:   group.Version,  // 群资料版本
			MemberVersion:  memberVersion,  // 群成员版本
			RequestVersion: requestVersion, // 入群申请版本
		}
	}

	// 转换为响应格式
	var result []types.GroupVersionItem
	for _, item := range groupVersions {
		result = append(result, *item)
	}

	resp = &types.GetUserGroupVersionsRes{
		Groups: result,
	}

	l.Infof("获取用户群组版本信息完成，用户ID: %s, 返回群组数: %d", req.UserID, len(result))
	return resp, nil
}
