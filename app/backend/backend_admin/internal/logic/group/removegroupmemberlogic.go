package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveGroupMemberLogic {
	return &RemoveGroupMemberLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RemoveGroupMemberLogic) RemoveGroupMember(req *types.RemoveGroupMemberReq) (resp *types.RemoveGroupMemberRes, err error) {
	if req.GroupId == "" {
		return nil, errors.New("群组ID不能为空")
	}
	if len(req.MemberIds) == 0 {
		return nil, errors.New("成员ID列表不能为空")
	}

	for _, userID := range req.MemberIds {
		if _, err = l.svcCtx.GroupRpc.RemoveGroupMember(l.ctx, &group_rpc.RemoveGroupMemberReq{
			GroupId: req.GroupId,
			UserId:  userID,
			Kick:    true,
		}); err != nil {
			l.Errorf("移除群成员失败 groupId=%s userId=%s: %v", req.GroupId, userID, err)
			return nil, err
		}
	}
	return &types.RemoveGroupMemberRes{}, nil
}
