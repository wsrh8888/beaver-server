package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
