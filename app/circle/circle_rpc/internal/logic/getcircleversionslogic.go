package logic

import (
	"context"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCircleVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleVersionsLogic {
	return &GetCircleVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCircleVersionsLogic) GetCircleVersions(in *circle_rpc.GetCircleVersionsReq) (*circle_rpc.GetCircleVersionsRes, error) {
	// 查用户加入的所有圈子
	var members []circle_models.CircleMemberModel
	l.svcCtx.DB.Where("user_id = ?", in.UserId).Find(&members)

	if len(members) == 0 {
		return &circle_rpc.GetCircleVersionsRes{List: []*circle_rpc.CircleVersionItem{}}, nil
	}

	circleIDs := make([]string, 0, len(members))
	roleMap := make(map[string]int32)
	for _, m := range members {
		circleIDs = append(circleIDs, m.CircleID)
		roleMap[m.CircleID] = int32(m.Role)
	}

	// 用 version 字段做增量比较（Redis 递增生成，严格单调递增）
	var circles []circle_models.CircleModel
	l.svcCtx.DB.Where("circle_id IN ? AND version > ?", circleIDs, in.Version).Find(&circles)

	list := make([]*circle_rpc.CircleVersionItem, 0, len(circles))
	for _, c := range circles {
		role := roleMap[c.CircleID]
		if c.IsDeleted {
			role = 0
		}
		list = append(list, &circle_rpc.CircleVersionItem{
			CircleId:    c.CircleID,
			Name:        c.Name,
			Avatar:      c.Avatar,
			MemberCount: c.MemberCount,
			Role:        role,
			Version:     c.Version,
		})
	}

	return &circle_rpc.GetCircleVersionsRes{List: list}, nil
}
