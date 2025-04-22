package logic

import (
	"context"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDisplayNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDisplayNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDisplayNameLogic {
	return &UpdateDisplayNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDisplayNameLogic) UpdateDisplayName(req *types.UpdateDisplayNameReq) (resp *types.UpdateDisplayNameRes, err error) {
	// todo: add your logic here and delete this line

	return
}
