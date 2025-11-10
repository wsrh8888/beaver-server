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

// 群资料同步
func NewGroupSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSyncLogic {
	return &GroupSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupSyncLogic) GroupSync(req *types.GroupSyncReq) (resp *types.GroupSyncRes, err error) {
	resp = &types.GroupSyncRes{
		Groups: []types.GroupSyncItem{},
	}

	if len(req.Groups) == 0 {
		l.Infof("群资料同步完成，用户ID: %s, 无需同步的群组", req.UserID)
		return resp, nil
	}

	// 为每个群组查询版本变化的数据
	for _, groupReq := range req.Groups {
		var groups []group_models.GroupModel
		err = l.svcCtx.DB.Where("group_id = ? AND version >= ?", groupReq.GroupID, groupReq.Version).
			Find(&groups).Error
		if err != nil {
			l.Errorf("查询群组数据失败，群组ID: %s, 错误: %v", groupReq.GroupID, err)
			continue
		}

		for _, group := range groups {
			// 判断群组是否被删除（通过状态字段）
			isDeleted := group.Status != 1 // 假设状态1为正常，其他为删除

			resp.Groups = append(resp.Groups, types.GroupSyncItem{
				GroupID:   group.GroupID,
				Title:     group.Title,
				Avatar:    group.Avatar,
				CreatorID: group.CreatorID,
				JoinType:  group.JoinType,
				IsDeleted: isDeleted,
				Version:   group.Version,
				CreateAt:  time.Time(group.CreatedAt).Unix(),
				UpdateAt:  time.Time(group.UpdatedAt).Unix(),
			})
		}
	}

	l.Infof("群资料同步完成，用户ID: %s, 返回群组变化数: %d", req.UserID, len(resp.Groups))
	return resp, nil
}
