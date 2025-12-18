package logic

import (
	"context"

	"beaver/app/mcp/mcp_api/internal/svc"
	"beaver/app/mcp/mcp_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListToolsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有可用工具列表
func NewListToolsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListToolsLogic {
	return &ListToolsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListToolsLogic) ListTools() (resp *types.ListToolsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
