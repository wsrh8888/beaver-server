package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除菜单
func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuLogic) DeleteMenu(req *types.DeleteMenuReq) (resp *types.DeleteMenuRes, err error) {
	// 检查菜单是否存在
	var menu backend_models.AdminSystemMenu
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&menu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("菜单不存在: %d", req.Id)
			return nil, err
		}
		logx.Errorf("查询菜单失败: %v", err)
		return nil, err
	}

	// 检查是否有子菜单
	var childCount int64
	err = l.svcCtx.DB.Model(&backend_models.AdminSystemMenu{}).Where("parent_id = ?", req.Id).Count(&childCount).Error
	if err != nil {
		logx.Errorf("检查子菜单失败: %v", err)
		return nil, err
	}

	if childCount > 0 {
		logx.Errorf("无法删除菜单，存在%d个子菜单", childCount)
		return nil, err
	}

	// 删除菜单
	err = l.svcCtx.DB.Delete(&menu).Error
	if err != nil {
		logx.Errorf("删除菜单失败: %v", err)
		return nil, err
	}

	// 删除相关的权限关联
	err = l.svcCtx.DB.Where("menu_id = ?", req.Id).Delete(&backend_models.AdminSystemAuthorityMenu{}).Error
	if err != nil {
		logx.Errorf("删除菜单权限关联失败: %v", err)
		// 这里不返回错误，因为菜单已经删除了
	}

	logx.Infof("菜单删除成功: ID=%d, Name=%s", req.Id, menu.Name)
	return &types.DeleteMenuRes{}, nil
}
