package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAgentAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAgentAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentAppLogic {
	return &DeleteAgentAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAgentAppLogic) DeleteAgentApp(in *platform_rpc.DeleteAgentAppReq) (*platform_rpc.DeleteAgentAppRes, error) {
	l.svcCtx.DB.Where("app_id = ?", in.AppId).Delete(&agent_models.AgentApp{})
	return &platform_rpc.DeleteAgentAppRes{}, nil
}
