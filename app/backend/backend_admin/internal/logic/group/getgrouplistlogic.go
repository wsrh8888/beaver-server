package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_models"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取群组列表
func NewGetGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupListLogic {
	return &GetGroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupListLogic) GetGroupList(req *types.GetGroupListReq) (resp *types.GetGroupListRes, err error) {
	// 分页参数校验
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("1 = 1")

	// 状态筛选
	if req.Status != 0 {
		whereClause = whereClause.Where("status = ?", req.Status)
	}

	// 类型筛选
	if req.Type != 0 {
		whereClause = whereClause.Where("type = ?", req.Type)
	}

	// 分页查询
	groups, count, err := list_query.ListQuery(l.svcCtx.DB, group_models.GroupModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  page,
			Limit: limit,
			Key:   req.Keywords,
			Sort:  "created_at desc",
		},
		Where: whereClause,
		Likes: []string{"title"},
	})

	if err != nil {
		logx.Errorf("查询群组列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.GetGroupListItem
	for _, group := range groups {
		list = append(list, types.GetGroupListItem{
			Id:        group.Id,
			Uuid:      group.GroupID,
			Type:      int(group.Type),
			Title:     group.Title,
			FileName:  group.Avatar,
			CreatorId: group.CreatorID,
			Notice:    group.Notice,
			Status:    int(group.Status),
			CreatedAt: group.CreatedAt.String(),
			UpdatedAt: group.UpdatedAt.String(),
		})
	}

	return &types.GetGroupListRes{
		List:  list,
		Total: count,
	}, nil
}
