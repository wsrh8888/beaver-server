package system

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAuthorityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAuthorityLogic {
	return &DeleteAuthorityLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteAuthorityLogic) DeleteAuthority(req *types.DeleteAuthorityReq) (resp *types.DeleteAuthorityRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("角色ID不能为空")
	}
	var userCount int64
	if err = l.svcCtx.DB.Model(&backend_models.AdminSystemAuthorityUser{}).
		Where("authority_id = ?", req.Id).Count(&userCount).Error; err != nil {
		return nil, err
	}
	if userCount > 0 {
		return nil, errors.New("该角色仍有关联管理员，无法删除")
	}
	if err = l.svcCtx.DB.Where("authority_id = ?", req.Id).
		Delete(&backend_models.AdminSystemAuthorityMenu{}).Error; err != nil {
		return nil, err
	}
	if err = l.svcCtx.DB.Delete(&backend_models.AdminSystemAuthority{}, req.Id).Error; err != nil {
		l.Errorf("删除角色失败: %v", err)
		return nil, err
	}
	return &types.DeleteAuthorityRes{}, nil
}
