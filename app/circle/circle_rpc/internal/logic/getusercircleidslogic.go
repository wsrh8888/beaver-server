package logic

import (
	"context"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserCircleIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserCircleIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCircleIDsLogic {
	return &GetUserCircleIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserCircleIDsLogic) GetUserCircleIDs(in *circle_rpc.GetUserCircleIDsReq) (*circle_rpc.GetUserCircleIDsRes, error) {
	var circleIDs []string
	l.svcCtx.DB.Model(&circle_models.CircleMemberModel{}).
		Where("user_id = ?", in.UserId).
		Pluck("circle_id", &circleIDs)

	return &circle_rpc.GetUserCircleIDsRes{CircleIds: circleIDs}, nil
}
