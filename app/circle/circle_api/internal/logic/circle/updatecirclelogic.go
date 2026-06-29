package circle

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCircleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCircleLogic {
	return &UpdateCircleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCircleLogic) UpdateCircle(req *types.UpdateCircleReq) (resp *types.UpdateCircleRes, err error) {
	// 权限校验：必须是圈主或管理员
	var member circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&member).Error; err != nil {
		return nil, fmt.Errorf("无权限")
	}
	if member.Role > 2 {
		return nil, fmt.Errorf("无权限，仅圈主和管理员可修改")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.JoinType != 0 {
		updates["join_type"] = req.JoinType
	}

	if len(updates) == 0 {
		return &types.UpdateCircleRes{}, nil
	}

	updates["version"] = l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", req.CircleID)
	if err = l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ? AND is_deleted = false", req.CircleID).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新圈子失败: %v", err)
	}

	return &types.UpdateCircleRes{}, nil
}
