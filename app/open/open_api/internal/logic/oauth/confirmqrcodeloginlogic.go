package oauth

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmQrCodeLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 确认扫码登录
func NewConfirmQrCodeLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmQrCodeLoginLogic {
	return &ConfirmQrCodeLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmQrCodeLoginLogic) ConfirmQrCodeLogin(req *types.ConfirmQrCodeLoginReq) (resp *types.ConfirmQrCodeLoginRes, err error) {
	// todo: add your logic here and delete this line

	return
}
