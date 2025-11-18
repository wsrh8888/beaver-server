package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAuthorityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建权限
func NewCreateAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAuthorityLogic {
	return &CreateAuthorityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAuthorityLogic) CreateAuthority(req *types.CreateAuthorityReq) (resp *types.CreateAuthorityRes, err error) {
	// todo: add your logic here and delete this line

	return
}
