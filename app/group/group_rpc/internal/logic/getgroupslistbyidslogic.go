package logic

import (
	"context"
	"time"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupsListByIdsLogic {
	return &GetGroupsListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupsListByIdsLogic) GetGroupsListByIds(in *group_rpc.GetGroupsListByIdsReq) (*group_rpc.GetGroupsListByIdsRes, error) {
	if len(in.GroupIDs) == 0 {
		return &group_rpc.GetGroupsListByIdsRes{Groups: []*group_rpc.GroupListById{}}, nil
	}

	// 查询指定群组ID列表中，自指定时间戳以来变更的群组资料
	var changedGroups []group_models.GroupModel
	query := l.svcCtx.DB.Where("group_id IN (?)", in.GroupIDs)
	if in.Since > 0 {
		query = query.Where("updated_at > ?", in.Since)
	}

	err := query.Find(&changedGroups).Error
	if err != nil {
		l.Errorf("查询变更的群组资料失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var groups []*group_rpc.GroupListById
	for _, group := range changedGroups {
		groups = append(groups, &group_rpc.GroupListById{
			GroupID:     group.GroupID,
			Name:        group.Title, // 使用Title字段作为群组名称
			Avatar:      group.Avatar,
			Description: group.Notice, // 使用Notice字段作为描述
			Version:     group.Version,
			CreatedAt:   time.Time(group.CreatedAt).UnixMilli(),
			UpdatedAt:   time.Time(group.UpdatedAt).UnixMilli(),
		})
	}

	return &group_rpc.GetGroupsListByIdsRes{Groups: groups}, nil
}
