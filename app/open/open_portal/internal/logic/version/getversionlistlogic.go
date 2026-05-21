package version

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVersionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取版本列表
func NewGetVersionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVersionListLogic {
	return &GetVersionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVersionListLogic) GetVersionList(req *types.GetVersionListReq) (resp *types.GetVersionListRes, err error) {
	return &types.GetVersionListRes{
		Total: 0,
		List:  []types.VersionInfo{},
	}, nil
}
