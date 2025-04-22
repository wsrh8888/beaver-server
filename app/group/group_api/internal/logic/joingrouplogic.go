package logic

import (
	"context"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinGroupLogic) JoinGroup(req *types.GroupJoinReq) (resp *types.GroupJoinRes, err error) {
	// todo: add your logic here and delete this line

	return
}
