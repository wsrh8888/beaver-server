package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MuteGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMuteGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteGroupMemberLogic {
	return &MuteGroupMemberLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *MuteGroupMemberLogic) MuteGroupMember(req *types.MuteGroupMemberReq) (resp *types.MuteGroupMemberRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("成员ID不能为空")
	}
	if req.ProhibitionTime < 0 {
		return nil, errors.New("禁言时长不能为负数")
	}

	minutes := int32(req.ProhibitionTime)
	_, err = l.svcCtx.GroupRpc.UpdateGroupMember(l.ctx, &group_rpc.UpdateGroupMemberReq{
		Id:          uint64(req.Id),
		MuteMinutes: &minutes,
	})
	if err != nil {
		l.Errorf("禁言群成员失败: %v", err)
		return nil, err
	}
	return &types.MuteGroupMemberRes{}, nil
}
