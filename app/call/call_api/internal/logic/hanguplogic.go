package logic

import (
	"context"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HangupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 主动挂断或拒绝通话
func NewHangupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HangupLogic {
	return &HangupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HangupLogic) Hangup(req *types.HangupCallReq) (resp *types.HangupCallRes, err error) {
	// todo: add your logic here and delete this line

	return
}
