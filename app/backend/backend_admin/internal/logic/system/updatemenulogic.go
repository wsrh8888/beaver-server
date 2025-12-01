package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新菜单
func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuReq) (resp *types.UpdateMenuRes, err error) {
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

	// 准备更新数据
	updateData := map[string]interface{}{
		"path":   req.Path,
		"name":   req.Name,
		"hidden": req.Hidden,
		"sort":   req.Sort,
		"title":  req.Title,
		"icon":   req.Icon,
	}

	// 处理parent_id
	if req.ParentId == 0 {
		updateData["parent_id"] = nil
	} else {
		updateData["parent_id"] = req.ParentId
	}

	// 更新菜单
	err = l.svcCtx.DB.Model(&menu).Updates(updateData).Error
	if err != nil {
		logx.Errorf("更新菜单失败: %v", err)
		return nil, err
	}

	logx.Infof("菜单更新成功: ID=%d, Name=%s", req.Id, req.Name)
	return &types.UpdateMenuRes{}, nil
}
