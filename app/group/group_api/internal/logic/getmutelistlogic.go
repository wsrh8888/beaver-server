package logic

import (
	"context"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMuteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMuteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMuteListLogic {
	return &GetMuteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMuteListLogic) GetMuteList(req *types.GroupMuteListReq) (resp *types.GroupMuteListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
