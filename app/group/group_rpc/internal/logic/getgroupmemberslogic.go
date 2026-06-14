package logic

import (
	"context"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMembersLogic {
	return &GetGroupMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMembersLogic) GetGroupMembers(in *group_rpc.GetGroupMembersReq) (*group_rpc.GetGroupMembersRes, error) {
	// 1. 查询群组成员
	var groupMembers []group_models.GroupMemberModel
	err := l.svcCtx.DB.Where("group_id = ?", in.GroupID).Find(&groupMembers).Error
	if err != nil {
		logx.Error("查询群组成员失败:", err)
		return nil, err
	}

	// 如果没有成员，返回空列表
	if len(groupMembers) == 0 {
		return &group_rpc.GetGroupMembersRes{
			Members: []*group_rpc.GroupMemberInfo{},
		}, nil
	}

	// 2. 收集所有成员的用户ID
	var userIDs []string
	for _, member := range groupMembers {
		userIDs = append(userIDs, member.UserID)
	}

	// 3. 通过 UserRpc 批量查询用户信息
	userResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
		UserIdList: userIDs,
	})
	if err != nil {
		logx.Error("查询用户信息失败:", err)
		return nil, err
	}

	// 4. 构建返回结果
	var members []*group_rpc.GroupMemberInfo
	for _, member := range groupMembers {
		if user, ok := userResp.UserInfo[member.UserID]; ok && user != nil {
			members = append(members, &group_rpc.GroupMemberInfo{
				UserID:   user.UserId,
				Username: user.NickName,
				Avatar:   user.Avatar,
			})
		}
	}

	return &group_rpc.GetGroupMembersRes{
		Members: members,
	}, nil
}
