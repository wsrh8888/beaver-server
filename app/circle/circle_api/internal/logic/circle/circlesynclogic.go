package circle

import (
	"context"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CircleSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCircleSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CircleSyncLogic {
	return &CircleSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CircleSyncLogic) CircleSync(req *types.CircleSyncReq) (resp *types.CircleSyncRes, err error) {
	var members []circle_models.CircleMemberModel
	l.svcCtx.DB.Where("user_id = ?", req.UserID).Find(&members)

	if len(members) == 0 {
		return &types.CircleSyncRes{List: []types.CircleSyncItem{}}, nil
	}

	circleIDs := make([]string, 0, len(members))
	roleMap := make(map[string]int8)
	for _, m := range members {
		circleIDs = append(circleIDs, m.CircleID)
		roleMap[m.CircleID] = m.Role
	}

	// 用 version 字段做增量比较
	var circles []circle_models.CircleModel
	l.svcCtx.DB.Where("circle_id IN ? AND version > ?", circleIDs, req.Version).Find(&circles)

	items := make([]types.CircleSyncItem, 0, len(circles))
	for _, c := range circles {
		role := roleMap[c.CircleID]
		if c.IsDeleted {
			role = 0
		}
		items = append(items, types.CircleSyncItem{
			CircleID:    c.CircleID,
			Name:        c.Name,
			Avatar:      c.Avatar,
			MemberCount: c.MemberCount,
			Role:        role,
			Version:     c.Version,
		})
	}

	return &types.CircleSyncRes{List: items}, nil
}
