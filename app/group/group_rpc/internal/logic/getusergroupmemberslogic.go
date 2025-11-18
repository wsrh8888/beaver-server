package logic

import (
	"context"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserGroupMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserGroupMembersLogic {
	return &GetUserGroupMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserGroupMembersLogic) GetUserGroupMembers(in *group_rpc.GetUserGroupMembersReq) (*group_rpc.GetUserGroupMembersRes, error) {
	// 获取用户加入的所有群组
	var userMembers []group_models.GroupMemberModel

	err := l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("user_id = ? AND status = ?", in.UserID, 1).
		Find(&userMembers).Error

	if err != nil {
		l.Errorf("查询用户群组失败: %v", err)
		return nil, err
	}

	// 提取用户加入的群组ID
	groupIDs := make([]string, 0, len(userMembers))
	for _, member := range userMembers {
		groupIDs = append(groupIDs, member.GroupID)
	}

	if len(groupIDs) == 0 {
		l.Infof("用户未加入任何群组，用户ID: %s", in.UserID)
		return &group_rpc.GetUserGroupMembersRes{
			MemberIDs: []string{},
		}, nil
	}

	// 获取所有群组的成员列表
	var allMembers []group_models.GroupMemberModel
	err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id IN ? AND status = ?", groupIDs, 1).
		Find(&allMembers).Error

	if err != nil {
		l.Errorf("查询群组成员失败: %v", err)
		return nil, err
	}

	// 去重并排除自己
	seen := make(map[string]bool)
	var allMemberIDs []string

	for _, member := range allMembers {
		if member.UserID != in.UserID && !seen[member.UserID] {
			seen[member.UserID] = true
			allMemberIDs = append(allMemberIDs, member.UserID)
		}
	}

	l.Infof("获取用户群成员成功，用户ID: %s, 群成员数: %d", in.UserID, len(allMemberIDs))

	return &group_rpc.GetUserGroupMembersRes{
		MemberIDs: allMemberIDs,
	}, nil
}
