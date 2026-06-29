package circle

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCircleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleDetailLogic {
	return &GetCircleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCircleDetailLogic) GetCircleDetail(req *types.GetCircleDetailReq) (resp *types.GetCircleDetailRes, err error) {
	var circle circle_models.CircleModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND is_deleted = false", req.CircleID).First(&circle).Error; err != nil {
		return nil, fmt.Errorf("圈子不存在")
	}

	role := int8(0)
	var member circle_models.CircleMemberModel
	if l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&member).Error == nil {
		role = member.Role
	}

	return &types.GetCircleDetailRes{
		CircleID:    circle.CircleID,
		Name:        circle.Name,
		Description: circle.Description,
		Avatar:      circle.Avatar,
		JoinType:    circle.JoinType,
		CreatorID:   circle.CreatorID,
		MemberCount: circle.MemberCount,
		PostCount:   circle.PostCount,
		Role:        role,
		CreatedAt:   circle.CreatedAt.String(),
	}, nil
}
