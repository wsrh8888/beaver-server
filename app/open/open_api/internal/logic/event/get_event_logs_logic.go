// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package event

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取事件推送日志
func NewGetEventLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventLogsLogic {
	return &GetEventLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventLogsLogic) GetEventLogs(req *types.GetEventLogsReq) (resp *types.GetEventLogsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
