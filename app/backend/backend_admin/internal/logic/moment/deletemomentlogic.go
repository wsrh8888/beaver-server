package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMomentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMomentLogic {
	return &DeleteMomentLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteMomentLogic) DeleteMoment(req *types.DeleteMomentReq) (resp *types.DeleteMomentRes, err error) {
	if req.MomentId == "" {
		return nil, errors.New("动态ID不能为空")
	}

	deleted := true
	_, err = l.svcCtx.MomentRpc.UpdateMoment(l.ctx, &moment_rpc.UpdateMomentReq{
		MomentId:  req.MomentId,
		IsDeleted: &deleted,
	})
	if err != nil {
		l.Errorf("删除动态失败: %v", err)
		return nil, err
	}
	return &types.DeleteMomentRes{}, nil
}
