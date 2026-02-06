package logic

import (
	"context"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartCallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发起音视频通话
func NewStartCallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartCallLogic {
	return &StartCallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartCallLogic) StartCall(req *types.StartCallReq) (resp *types.StartCallRes, err error) {
	// todo: add your logic here and delete this line

	return
}
