package robot

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveRobotFromGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 将 Robot 从群移除
func NewRemoveRobotFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveRobotFromGroupLogic {
	return &RemoveRobotFromGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveRobotFromGroupLogic) RemoveRobotFromGroup(req *types.RemoveRobotFromGroupReq) (resp *types.RemoveRobotFromGroupRes, err error) {
	// todo: add your logic here and delete this line

	return
}
