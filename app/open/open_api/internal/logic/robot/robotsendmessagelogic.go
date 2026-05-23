package robot

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RobotSendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Robot 发送消息到 IM 会话（私聊或群聊）
func NewRobotSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RobotSendMessageLogic {
	return &RobotSendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RobotSendMessageLogic) RobotSendMessage(req *types.RobotSendMessageReq) (resp *types.RobotSendMessageRes, err error) {
	// todo: add your logic here and delete this line

	return
}
