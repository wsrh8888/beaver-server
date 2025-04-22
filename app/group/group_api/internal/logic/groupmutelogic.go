package logic

import (
	"context"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMuteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMuteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMuteLogic {
	return &GroupMuteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMuteLogic) GroupMute(req *types.GroupMuteReq) (resp *types.GroupMuteRes, err error) {
	// todo: add your logic here and delete this line

	return
}
