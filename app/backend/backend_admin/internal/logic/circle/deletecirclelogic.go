package circle

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCircleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCircleLogic {
	return &DeleteCircleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCircleLogic) DeleteCircle(req *types.DeleteCircleReq) (resp *types.DeleteCircleRes, err error) {
	if req.CircleId == "" {
		return nil, errors.New("圈子ID不能为空")
	}

	deleted := true
	_, err = l.svcCtx.CircleRpc.UpdateCircle(l.ctx, &circle_rpc.UpdateCircleReq{
		CircleId:  req.CircleId,
		IsDeleted: &deleted,
	})
	if err != nil {
		l.Errorf("解散圈子失败: %v", err)
		return nil, err
	}
	return &types.DeleteCircleRes{}, nil
}
