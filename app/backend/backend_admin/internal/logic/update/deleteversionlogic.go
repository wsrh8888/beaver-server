package update

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVersionLogic {
	return &DeleteVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteVersionLogic) DeleteVersion(req *types.DeleteVersionReq) (resp *types.DeleteVersionRes, err error) {
	_, err = l.svcCtx.PlatformRpc.DeleteVersion(l.ctx, &platform_rpc.DeleteVersionReq{
		VersionId: uint64(req.VersionID),
	})
	if err != nil {
		l.Errorf("删除版本失败: %v", err)
		return nil, err
	}
	return &types.DeleteVersionRes{}, nil
}
