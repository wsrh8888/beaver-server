package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取菜单列表
func NewGetMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenusLogic {
	return &GetMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenusLogic) GetMenus(req *types.GetMenuListReq) (resp *types.GetMenuListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
