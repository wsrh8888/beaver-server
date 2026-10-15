package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAgentModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentModelLogic {
	return &DeleteAgentModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAgentModelLogic) DeleteAgentModel(in *platform_rpc.DeleteAgentModelReq) (*platform_rpc.DeleteAgentModelRes, error) {
	q := l.svcCtx.DB.Where("model_id = ?", in.ModelId)
	if in.Source != "" {
		q = q.Where("source = ?", in.Source)
	}
	if in.UserId != "" {
		q = q.Where("user_id = ?", in.UserId)
	}
	q.Delete(&agent_models.AgentModel{})
	return &platform_rpc.DeleteAgentModelRes{}, nil
}
