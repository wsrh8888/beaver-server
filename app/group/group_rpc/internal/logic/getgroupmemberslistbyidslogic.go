package logic

import (
	"context"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMembersListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMembersListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMembersListByIdsLogic {
	return &GetGroupMembersListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMembersListByIdsLogic) GetGroupMembersListByIds(in *group_rpc.GetGroupMembersListByIdsReq) (*group_rpc.GetGroupMembersListByIdsRes, error) {
	if len(in.GroupIDs) == 0 {
		return &group_rpc.GetGroupMembersListByIdsRes{Members: []*group_rpc.GroupMemberListById{}}, nil
	}

	// 查询指定群组ID列表中，自指定时间戳以来变更的群成员
	var changedMembers []group_models.GroupMemberModel
	query := l.svcCtx.DB.Where("group_id IN (?)", in.GroupIDs)
	if in.Since > 0 {
		query = query.Where("updated_at > ?", in.Since)
	}

	err := query.Find(&changedMembers).Error
	if err != nil {
		l.Errorf("查询变更的群成员失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var members []*group_rpc.GroupMemberListById
	for _, member := range changedMembers {
		members = append(members, &group_rpc.GroupMemberListById{
			GroupID:  member.GroupID,
			UserID:   member.UserID,
			Role:     int32(member.Role),
			JoinedAt: member.JoinTime.UnixMilli(),
			Version:  member.Version,
		})
	}

	return &group_rpc.GetGroupMembersListByIdsRes{Members: members}, nil
}
