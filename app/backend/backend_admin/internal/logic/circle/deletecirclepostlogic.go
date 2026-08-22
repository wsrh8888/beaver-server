package circle

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCirclePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCirclePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCirclePostLogic {
	return &DeleteCirclePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCirclePostLogic) DeleteCirclePost(req *types.DeleteCirclePostReq) (resp *types.DeleteCirclePostRes, err error) {
	if req.PostId == "" {
		return nil, errors.New("帖子ID不能为空")
	}

	_, err = l.svcCtx.CircleRpc.DeletePost(l.ctx, &circle_rpc.DeletePostReq{
		PostId: req.PostId,
	})
	if err != nil {
		l.Errorf("删除圈子帖子失败: %v", err)
		return nil, err
	}
	return &types.DeleteCirclePostRes{}, nil
}
