package logic

import (
	"context"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群组数据同步
func NewGroupSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSyncLogic {
	return &GroupSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupSyncLogic) GroupSync(req *types.GroupSyncReq) (resp *types.GroupSyncRes, err error) {
	var groupMembers []group_models.GroupMemberModel

	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	// 查询用户加入的群组成员关系
	err = l.svcCtx.DB.Where("user_id = ? AND version > ? AND version <= ?",
		req.UserID, req.FromVersion, req.ToVersion).
		Order("version ASC").
		Limit(limit + 1).
		Find(&groupMembers).Error
	if err != nil {
		l.Errorf("查询群组成员数据失败: %v", err)
		return nil, err
	}

	// 判断是否还有更多数据
	hasMore := len(groupMembers) > limit
	if hasMore {
		groupMembers = groupMembers[:limit]
	}

	// 转换为响应格式
	var groupItems []types.GroupSyncItem
	var nextVersion int64 = req.FromVersion

	for _, member := range groupMembers {
		// 查询群组信息
		var group group_models.GroupModel
		err = l.svcCtx.DB.Where("uuid = ?", member.GroupID).First(&group).Error
		if err != nil {
			l.Errorf("查询群组信息失败，群组ID: %s, 错误: %v", member.GroupID, err)
			continue
		}

		// 判断群组是否被删除（通过状态字段）
		isDeleted := group.Status != 1 // 假设状态1为正常，其他为删除

		groupItems = append(groupItems, types.GroupSyncItem{
			GroupID:   group.GroupID,
			Title:     group.Title,
			Avatar:    group.Avatar,
			CreatorID: group.CreatorID,
			JoinType:  group.JoinType,
			IsDeleted: isDeleted,
			Version:   member.Version, // 使用成员关系的版本号
			CreateAt:  time.Time(member.CreatedAt).Unix(),
			UpdateAt:  time.Time(member.UpdatedAt).Unix(),
		})

		nextVersion = member.Version
	}

	// 如果没有更多数据，nextVersion应该是toVersion+1
	if !hasMore {
		nextVersion = req.ToVersion + 1
	}

	resp = &types.GroupSyncRes{
		Groups:      groupItems,
		HasMore:     hasMore,
		NextVersion: nextVersion,
	}

	l.Infof("群组数据同步完成，用户ID: %s, 返回群组数: %d, 还有更多: %v", req.UserID, len(groupItems), hasMore)
	return resp, nil
}
