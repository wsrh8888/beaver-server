package logic

import (
	"context"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群成员数据同步
func NewGroupMemberSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberSyncLogic {
	return &GroupMemberSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberSyncLogic) GroupMemberSync(req *types.GroupMemberSyncReq) (resp *types.GroupMemberSyncRes, err error) {
	var groupMembers []group_models.GroupMemberModel

	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	// 查询用户作为成员的所有群成员记录
	err = l.svcCtx.DB.Where("user_id = ? AND version > ? AND version <= ?",
		req.UserID, req.FromVersion, req.ToVersion).
		Order("version ASC").
		Limit(limit + 1).
		Find(&groupMembers).Error
	if err != nil {
		l.Errorf("查询群成员数据失败: %v", err)
		return nil, err
	}

	// 判断是否还有更多数据
	hasMore := len(groupMembers) > limit
	if hasMore {
		groupMembers = groupMembers[:limit]
	}

	// 转换为响应格式
	var groupMemberItems []types.GroupMemberSyncItem
	var nextVersion int64 = req.FromVersion

	for _, member := range groupMembers {
		groupMemberItems = append(groupMemberItems, types.GroupMemberSyncItem{
			GroupID:  member.GroupID,
			UserID:   member.UserID,
			Role:     member.Role,
			Status:   member.Status,
			JoinTime: time.Time(member.JoinTime).Unix(),
			Version:  member.Version,
			CreateAt: time.Time(member.CreatedAt).Unix(),
			UpdateAt: time.Time(member.UpdatedAt).Unix(),
		})

		nextVersion = member.Version
	}

	// 如果没有更多数据，nextVersion应该是toVersion+1
	if !hasMore {
		nextVersion = req.ToVersion + 1
	}

	resp = &types.GroupMemberSyncRes{
		GroupMembers: groupMemberItems,
		HasMore:      hasMore,
		NextVersion:  nextVersion,
	}

	l.Infof("群成员数据同步完成，用户ID: %s, 返回群成员记录数: %d, 还有更多: %v", req.UserID, len(groupMemberItems), hasMore)
	return resp, nil
}
