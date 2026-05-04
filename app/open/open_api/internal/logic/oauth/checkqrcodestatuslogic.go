package oauth

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckQrCodeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询扫码状态
func NewCheckQrCodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckQrCodeStatusLogic {
	return &CheckQrCodeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckQrCodeStatusLogic) CheckQrCodeStatus(req *types.CheckQrCodeStatusReq) (resp *types.CheckQrCodeStatusRes, err error) {
	// todo: add your logic here and delete this line

	return
}
