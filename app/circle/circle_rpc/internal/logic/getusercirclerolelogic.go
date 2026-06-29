package logic

import (
	"context"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserCircleRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserCircleRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCircleRoleLogic {
	return &GetUserCircleRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserCircleRoleLogic) GetUserCircleRole(in *circle_rpc.GetUserCircleRoleReq) (*circle_rpc.GetUserCircleRoleRes, error) {
	var member circle_models.CircleMemberModel
	if err := l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", in.CircleId, in.UserId).
		First(&member).Error; err != nil {
		// 未加入返回 0
		return &circle_rpc.GetUserCircleRoleRes{Role: 0}, nil
	}

	return &circle_rpc.GetUserCircleRoleRes{Role: int32(member.Role)}, nil
}
