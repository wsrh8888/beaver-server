package digital

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	beaverlog "beaver/utils/beaverlog"
)

// goctl 会把 logx 写进 imports，占位避免未使用导入。
var _ = logx.WithContext

type ListDigitalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 数字伙伴列表
func NewListDigitalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDigitalLogic {
	return &ListDigitalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListDigital", ctx),
	}
}

func (l *ListDigitalLogic) ListDigital(req *types.ListDigitalReq) (resp *types.ListDigitalRes, err error) {
	// todo: add your logic here and delete this line

	return
}
