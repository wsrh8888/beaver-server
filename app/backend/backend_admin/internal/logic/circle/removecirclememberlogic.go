package circle

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveCircleMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveCircleMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveCircleMemberLogic {
	return &RemoveCircleMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveCircleMemberLogic) RemoveCircleMember(req *types.RemoveCircleMemberReq) (resp *types.RemoveCircleMemberRes, err error) {
	if req.CircleId == "" {
		return nil, errors.New("圈子ID不能为空")
	}
	if len(req.MemberIds) == 0 {
		return nil, errors.New("成员ID不能为空")
	}

	_, err = l.svcCtx.CircleRpc.RemoveCircleMembers(l.ctx, &circle_rpc.RemoveCircleMembersReq{
		CircleId: req.CircleId,
		UserIds:  req.MemberIds,
	})
	if err != nil {
		l.Errorf("移除圈子成员失败: %v", err)
		return nil, err
	}
	return &types.RemoveCircleMemberRes{}, nil
}
