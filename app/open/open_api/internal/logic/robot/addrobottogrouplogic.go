package robot

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRobotToGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 将 Robot 加入群（Robot 进群后可接收消息和 @ 提及，Beaver 会推送事件到 Robot 的 Webhook URL）
func NewAddRobotToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRobotToGroupLogic {
	return &AddRobotToGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRobotToGroupLogic) AddRobotToGroup(req *types.AddRobotToGroupReq) (resp *types.AddRobotToGroupRes, err error) {
	// todo: add your logic here and delete this line

	return
}
