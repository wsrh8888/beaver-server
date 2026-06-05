package system

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAuthorityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAuthorityLogic {
	return &UpdateAuthorityLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateAuthorityLogic) UpdateAuthority(req *types.UpdateAuthorityReq) (resp *types.UpdateAuthorityRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("角色ID不能为空")
	}
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"status":      req.Status,
		"sort":        req.Sort,
	}
	if err = l.svcCtx.DB.Model(&backend_models.AdminSystemAuthority{}).
		Where("id = ?", req.Id).Updates(updates).Error; err != nil {
		l.Errorf("更新角色失败: %v", err)
		return nil, err
	}
	return &types.UpdateAuthorityRes{}, nil
}
