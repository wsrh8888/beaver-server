package logic

import (
	"context"

	"beaver/app/mcp/mcp_api/internal/svc"
	"beaver/app/mcp/mcp_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteToolLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 执行MCP工具
func NewExecuteToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteToolLogic {
	return &ExecuteToolLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExecuteToolLogic) ExecuteTool(req *types.ExecuteToolReq) (resp *types.ExecuteToolRes, err error) {
	// todo: add your logic here and delete this line

	return
}
