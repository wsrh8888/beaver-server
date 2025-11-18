package logic

import (
	"context"
	"strings"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchGroupsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索群组
func NewSearchGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchGroupsLogic {
	return &SearchGroupsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchGroupsLogic) SearchGroups(req *types.GroupSearchReq) (resp *types.GroupSearchRes, err error) {
	resp = &types.GroupSearchRes{
		List: []types.GroupSearchItem{},
	}

	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // 限制最大每页数量
	}
	offset := (page - 1) * limit

	// 构建查询条件
	query := l.svcCtx.DB.Model(&group_models.GroupModel{}).Where("status = 1") // 只搜索正常状态的群组

	// 如果有搜索关键词，按群组名称模糊搜索（大小写不敏感）
	if req.Keyword != "" {
		// 使用ILIKE进行大小写不敏感的搜索，如果数据库不支持，则回退到LIKE
		keyword := strings.TrimSpace(req.Keyword)
		if keyword != "" {
			// 尝试使用ILIKE（PostgreSQL风格的大小写不敏感搜索）
			query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+keyword+"%")
		}
	}

	// 获取总数
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		l.Errorf("查询群组总数失败: %v", err)
		return nil, err
	}

	// 分页查询群组信息
	var groups []group_models.GroupModel
	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&groups).Error
	if err != nil {
		l.Errorf("查询群组列表失败: %v", err)
		return nil, err
	}

	// 为每个群组获取成员数量
	for _, group := range groups {
		var memberCount int64
		err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
			Where("group_id = ? AND status = 1", group.GroupID).
			Count(&memberCount).Error
		if err != nil {
			l.Errorf("查询群组成员数量失败, groupId: %s, error: %v", group.GroupID, err)
			memberCount = 0
		}

		resp.List = append(resp.List, types.GroupSearchItem{
			GroupID:     group.GroupID,
			Title:       group.Title,
			Avatar:      group.Avatar,
			MemberCount: int(memberCount),
			JoinType:    group.JoinType,
			CreatorID:   group.CreatorID,
		})
	}

	resp.Count = total

	l.Infof("搜索群组完成，关键词: %s, 返回群组数: %d, 总数: %d", req.Keyword, len(resp.List), total)
	return resp, nil
}
