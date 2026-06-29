package logic

import (
	"context"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCircleMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleMembersLogic {
	return &GetCircleMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCircleMembersLogic) GetCircleMembers(in *circle_rpc.GetCircleMembersReq) (*circle_rpc.GetCircleMembersRes, error) {
	var total int64
	var members []circle_models.CircleMemberModel

	l.svcCtx.DB.Model(&circle_models.CircleMemberModel{}).
		Where("circle_id = ?", in.CircleId).
		Count(&total)
	l.svcCtx.DB.Where("circle_id = ?", in.CircleId).
		Offset(int((in.Page - 1) * in.PageSize)).
		Limit(int(in.PageSize)).
		Find(&members)

	list := make([]*circle_rpc.CircleMemberItem, 0, len(members))
	for _, m := range members {
		list = append(list, &circle_rpc.CircleMemberItem{
			CircleId: m.CircleID,
			UserId:   m.UserID,
			Role:     int32(m.Role),
		})
	}

	return &circle_rpc.GetCircleMembersRes{Total: total, List: list}, nil
}
