package event

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TestEventPushLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 触发测试事件推送（调试用，向 Bot 服务器发一条测试事件）
func NewTestEventPushLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestEventPushLogic {
	return &TestEventPushLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestEventPushLogic) TestEventPush(req *types.TestEventPushReq) (resp *types.TestEventPushRes, err error) {
	// todo: add your logic here and delete this line

	return
}
