package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAuthorityMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新权限菜单
func NewUpdateAuthorityMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAuthorityMenuLogic {
	return &UpdateAuthorityMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAuthorityMenuLogic) UpdateAuthorityMenu(req *types.UpdateAuthorityMenuReq) (resp *types.UpdateAuthorityMenuRes, err error) {
	// 先删除该角色原有的所有菜单关联
	err = l.svcCtx.DB.Where("authority_id = ?", req.Id).Delete(&backend_models.AdminSystemAuthorityMenu{}).Error
	if err != nil {
		logx.Errorf("删除原有的菜单关联失败: %v", err)
		return nil, err
	}

	// 如果有新的菜单权限，则批量插入
	if len(req.Menus) > 0 {
		var authorityMenus []backend_models.AdminSystemAuthorityMenu
		for _, menu := range req.Menus {
			authorityMenus = append(authorityMenus, backend_models.AdminSystemAuthorityMenu{
				AuthorityID: req.Id,
				MenuID:      menu.Id,
			})
		}

		err = l.svcCtx.DB.Create(&authorityMenus).Error
		if err != nil {
			logx.Errorf("创建新的菜单关联失败: %v", err)
			return nil, err
		}
	}

	logx.Infof("权限菜单更新成功: 权限ID=%d, 菜单数量=%d", req.Id, len(req.Menus))
	return &types.UpdateAuthorityMenuRes{}, nil
}
