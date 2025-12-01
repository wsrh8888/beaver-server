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

type GetGroupMemberListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取群成员列表
func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupMemberListLogic) GetGroupMemberList(req *types.GetGroupMemberListReq) (resp *types.GetGroupMemberListRes, err error) {
	// 分页参数校验
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("group_id = ?", req.GroupId)

	// 角色筛选
	if req.Role != 0 {
		whereClause = whereClause.Where("role = ?", req.Role)
	}

	// 状态筛选
	if req.Status != 0 {
		whereClause = whereClause.Where("status = ?", req.Status)
	}

	// 分页查询
	members, count, err := list_query.ListQuery(l.svcCtx.DB, group_models.GroupMemberModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  page,
			Limit: limit,
			Sort:  "created_at desc",
		},
		Where: whereClause,
	})

	if err != nil {
		logx.Errorf("查询群组成员列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式（已精简模型字段，以下为兼容管理端返回结构的占位）
	var list []types.GetGroupMemberListItem
	for _, member := range members {
		list = append(list, types.GetGroupMemberListItem{
			Id:              member.Id,
			GroupId:         member.GroupID,
			UserId:          member.UserID,
			MemberNickname:  "",
			Role:            int(member.Role),
			ProhibitionTime: 0,
			InviterId:       "",
			Status:          int(member.Status),
			NotifyLevel:     0,
			DisplayName:     "",
			CreatedAt:       member.CreatedAt.String(),
			UpdatedAt:       member.UpdatedAt.String(),
		})
	}

	return &types.GetGroupMemberListRes{
		List:  list,
		Total: count,
	}, nil
}
