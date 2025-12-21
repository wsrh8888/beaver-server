package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建菜单
func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuReq) (resp *types.CreateMenuRes, err error) {
	// 创建菜单数据
	menu := backend_models.AdminSystemMenu{
		Path:   req.Path,
		Name:   req.Name,
		Hidden: req.Hidden,
		Sort:   req.Sort,
		Title:  req.Title,
		Icon:   req.Icon,
		Status: 1, // 默认启用
	}

	// 处理parent_id
	if req.ParentId != 0 {
		menu.ParentID = &req.ParentId
	}

	// 创建菜单
	err = l.svcCtx.DB.Create(&menu).Error
	if err != nil {
		logx.Errorf("创建菜单失败: %v", err)
		return nil, err
	}

	logx.Infof("菜单创建成功: ID=%d, Name=%s", menu.Id, menu.Name)
	return &types.CreateMenuRes{}, nil
}
