package circle

import (
	"context"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyCircleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyCircleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyCircleListLogic {
	return &GetMyCircleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyCircleListLogic) GetMyCircleList(req *types.GetMyCircleListReq) (resp *types.GetMyCircleListRes, err error) {
	var members []circle_models.CircleMemberModel
	var total int64

	l.svcCtx.DB.Model(&circle_models.CircleMemberModel{}).
		Where("user_id = ?", req.UserID).
		Count(&total)
	l.svcCtx.DB.Where("user_id = ?", req.UserID).
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&members)

	if len(members) == 0 {
		return &types.GetMyCircleListRes{Count: total, List: []types.MyCircleListItem{}}, nil
	}

	circleIDs := make([]string, 0, len(members))
	roleMap := make(map[string]int8)
	for _, m := range members {
		circleIDs = append(circleIDs, m.CircleID)
		roleMap[m.CircleID] = m.Role
	}

	var circles []circle_models.CircleModel
	l.svcCtx.DB.Where("circle_id IN ? AND is_deleted = false", circleIDs).Find(&circles)

	items := make([]types.MyCircleListItem, 0, len(circles))
	for _, c := range circles {
		items = append(items, types.MyCircleListItem{
			CircleID:    c.CircleID,
			Name:        c.Name,
			Avatar:      c.Avatar,
			MemberCount: c.MemberCount,
			PostCount:   c.PostCount,
			Role:        roleMap[c.CircleID],
		})
	}

	return &types.GetMyCircleListRes{Count: total, List: items}, nil
}
