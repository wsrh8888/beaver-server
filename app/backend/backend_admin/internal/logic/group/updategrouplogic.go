package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupLogic {
	return &UpdateGroupLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateGroupLogic) UpdateGroup(req *types.UpdateGroupReq) (resp *types.UpdateGroupRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("群组ID不能为空")
	}

	rpcReq := &group_rpc.UpdateGroupReq{
		Id:     uint64(req.Id),
		Title:  req.Title,
		Avatar: req.FileName,
		Notice: req.Notice,
		Status: int32(req.Status),
	}
	muteAll := req.MuteAll
	rpcReq.MuteAll = &muteAll

	_, err = l.svcCtx.GroupRpc.UpdateGroup(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("更新群组信息失败: %v", err)
		return nil, err
	}
	return &types.UpdateGroupRes{}, nil
}
