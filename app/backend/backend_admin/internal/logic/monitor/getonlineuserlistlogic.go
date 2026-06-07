package monitor

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOnlineUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 在线用户列表
func NewGetOnlineUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOnlineUserListLogic {
	return &GetOnlineUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOnlineUserListLogic) GetOnlineUserList(req *types.GetOnlineUserListReq) (resp *types.GetOnlineUserListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
