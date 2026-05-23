package robot

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RobotStreamMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Robot 流式发送消息（SSE，适合 AI 流式输出）
func NewRobotStreamMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RobotStreamMessageLogic {
	return &RobotStreamMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RobotStreamMessageLogic) RobotStreamMessage(req *types.RobotStreamMessageReq) (resp *types.RobotStreamChunk, err error) {
	// todo: add your logic here and delete this line

	return
}
