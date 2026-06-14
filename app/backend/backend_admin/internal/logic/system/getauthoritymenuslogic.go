package system

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuthorityMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAuthorityMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorityMenusLogic {
	return &GetAuthorityMenusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAuthorityMenusLogic) GetAuthorityMenus(req *types.GetAuthorityMenusReq) (resp *types.GetAuthorityMenusRes, err error) {
	var rows []backend_models.AdminSystemAuthorityMenu
	if err = l.svcCtx.DB.Where("authority_id = ?", req.Id).Find(&rows).Error; err != nil {
		l.Errorf("查询角色菜单失败: %v", err)
		return nil, err
	}
	menuIds := make([]uint, 0, len(rows))
	for _, row := range rows {
		menuIds = append(menuIds, uint(row.MenuID))
	}
	return &types.GetAuthorityMenusRes{MenuIds: menuIds}, nil
}
