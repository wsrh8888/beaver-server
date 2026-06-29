package logic

import (
	"context"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCircleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleListLogic {
	return &GetCircleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCircleListLogic) GetCircleList(in *circle_rpc.GetCircleListReq) (*circle_rpc.GetCircleListRes, error) {
	query := l.svcCtx.DB.Model(&circle_models.CircleModel{}).Where("is_deleted = false")

	if in.CircleId != "" {
		query = query.Where("circle_id = ?", in.CircleId)
	}
	if in.UserId != "" {
		// 查询用户加入的圈子ID
		var memberCircleIDs []string
		l.svcCtx.DB.Model(&circle_models.CircleMemberModel{}).
			Where("user_id = ?", in.UserId).
			Pluck("circle_id", &memberCircleIDs)
		query = query.Where("circle_id IN ?", memberCircleIDs)
	}
	if in.Keywords != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+in.Keywords+"%", "%"+in.Keywords+"%")
	}

	var total int64
	query.Count(&total)

	var circles []circle_models.CircleModel
	query.Order("created_at DESC").
		Offset(int((in.Page - 1) * in.PageSize)).
		Limit(int(in.PageSize)).
		Find(&circles)

	list := make([]*circle_rpc.CircleItem, 0, len(circles))
	for _, c := range circles {
		list = append(list, &circle_rpc.CircleItem{
			CircleId:    c.CircleID,
			Name:        c.Name,
			Description: c.Description,
			Avatar:      c.Avatar,
			CreatorId:   c.CreatorID,
			JoinType:    int32(c.JoinType),
			MemberCount: c.MemberCount,
			PostCount:   c.PostCount,
			IsDeleted:   c.IsDeleted,
			CreatedAt:   c.CreatedAt.String(),
			UpdatedAt:   c.UpdatedAt.String(),
		})
	}

	return &circle_rpc.GetCircleListRes{Total: total, List: list}, nil
}
