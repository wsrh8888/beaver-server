package logic

import (
	"context"

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
	var groups []struct {
		GroupID string `gorm:"column:group_id"`
	}

	err := l.svcCtx.DB.Raw(`
		SELECT group_id
		FROM group_members
		WHERE user_id = ? AND status = 1
	`, in.UserID).Scan(&groups).Error

	if err != nil {
		l.Errorf("查询用户群组失败: %v", err)
		return nil, err
	}

	var allMemberIDs []string
	seen := make(map[string]bool)

	// 遍历每个群组，获取成员列表
	for _, group := range groups {
		// 获取群成员列表
		var members []struct {
			UserID string `gorm:"column:user_id"`
		}

		err := l.svcCtx.DB.Raw(`
			SELECT user_id
			FROM group_members
			WHERE group_id = ? AND status = 1
		`, group.GroupID).Scan(&members).Error

		if err != nil {
			l.Errorf("查询群组成员失败，群组ID: %s, 错误: %v", group.GroupID, err)
			continue
		}

		// 添加群成员ID（排除自己，并且去重）
		for _, member := range members {
			if member.UserID != in.UserID && !seen[member.UserID] {
				seen[member.UserID] = true
				allMemberIDs = append(allMemberIDs, member.UserID)
			}
		}
	}

	l.Infof("获取用户群成员成功，用户ID: %s, 群成员数: %d", in.UserID, len(allMemberIDs))

	return &group_rpc.GetUserGroupMembersRes{
		MemberIDs: allMemberIDs,
	}, nil
}
