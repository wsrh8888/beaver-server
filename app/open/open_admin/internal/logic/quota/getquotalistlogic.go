package quota

import (
	"context"

	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuotaListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQuotaListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuotaListLogic {
	return &GetQuotaListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQuotaListLogic) GetQuotaList(req *types.GetQuotaListReq) (resp *types.GetQuotaListRes, err error) {
	// TODO: 实现获取配额列表逻辑
	return &types.GetQuotaListRes{
		Total: 0,
		List:  []types.QuotaInfo{},
	}, nil
}
