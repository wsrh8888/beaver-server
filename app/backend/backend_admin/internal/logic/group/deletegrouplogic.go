package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGroupLogic {
	return &DeleteGroupLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteGroupLogic) DeleteGroup(req *types.DeleteGroupReq) (resp *types.DeleteGroupRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("群组ID不能为空")
	}

	_, err = l.svcCtx.GroupRpc.UpdateGroup(l.ctx, &group_rpc.UpdateGroupReq{
		Id:     uint64(req.Id),
		Status: 3,
	})
	if err != nil {
		l.Errorf("解散群组失败: %v", err)
		return nil, err
	}
	return &types.DeleteGroupRes{}, nil
}
