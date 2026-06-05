package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMemberRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMemberRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMemberRoleLogic {
	return &UpdateMemberRoleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateMemberRoleLogic) UpdateMemberRole(req *types.UpdateMemberRoleReq) (resp *types.UpdateMemberRoleRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("成员ID不能为空")
	}
	if req.Role < 1 || req.Role > 3 {
		return nil, errors.New("角色值无效，有效值为1-3")
	}

	role := int32(req.Role)
	_, err = l.svcCtx.GroupRpc.UpdateGroupMember(l.ctx, &group_rpc.UpdateGroupMemberReq{
		Id:   uint64(req.Id),
		Role: &role,
	})
	if err != nil {
		l.Errorf("更新成员角色失败: %v", err)
		return nil, err
	}
	return &types.UpdateMemberRoleRes{}, nil
}
