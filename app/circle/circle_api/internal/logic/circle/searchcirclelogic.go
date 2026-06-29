package circle

import (
	"context"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCircleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCircleLogic {
	return &SearchCircleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchCircleLogic) SearchCircle(req *types.SearchCircleReq) (resp *types.SearchCircleRes, err error) {
	var total int64
	var circles []circle_models.CircleModel

	query := l.svcCtx.DB.Model(&circle_models.CircleModel{}).Where("is_deleted = false")
	if req.Keywords != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+req.Keywords+"%", "%"+req.Keywords+"%")
	}
	query.Count(&total)
	query.Order("member_count DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&circles)

	if len(circles) == 0 {
		return &types.SearchCircleRes{Count: total, List: []types.SearchCircleItem{}}, nil
	}

	// 批量查当前用户在这些圈子的角色
	circleIDs := make([]string, 0, len(circles))
	for _, c := range circles {
		circleIDs = append(circleIDs, c.CircleID)
	}
	var members []circle_models.CircleMemberModel
	l.svcCtx.DB.Where("circle_id IN ? AND user_id = ?", circleIDs, req.UserID).Find(&members)
	roleMap := make(map[string]int8)
	for _, m := range members {
		roleMap[m.CircleID] = m.Role
	}

	items := make([]types.SearchCircleItem, 0, len(circles))
	for _, c := range circles {
		items = append(items, types.SearchCircleItem{
			CircleID:    c.CircleID,
			Name:        c.Name,
			Description: c.Description,
			Avatar:      c.Avatar,
			MemberCount: c.MemberCount,
			JoinType:    c.JoinType,
			Role:        roleMap[c.CircleID],
		})
	}

	return &types.SearchCircleRes{Count: total, List: items}, nil
}
