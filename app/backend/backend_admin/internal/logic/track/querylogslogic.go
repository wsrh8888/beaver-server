package logic

import (
	"context"
	"encoding/json"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryLogsLogic {
	return &QueryLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryLogsLogic) QueryLogs(req *types.QueryLogsReq) (resp *types.QueryLogsRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.AdminQueryLogs(l.ctx, &platform_rpc.AdminQueryLogsReq{
		BucketId:   req.BucketID,
		Level:      req.Level,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Keyword:    req.Keyword,
		UserFilter: req.UserFilter,
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
	})
	if err != nil {
		l.Errorf("查询日志失败: %v", err)
		return nil, err
	}

	logs := make([]types.QueryLogsItem, 0, len(rpcRes.Logs))
	for _, item := range rpcRes.Logs {
		var data interface{}
		if item.Data != "" {
			if err := json.Unmarshal([]byte(item.Data), &data); err != nil {
				data = item.Data
			}
		}
		logs = append(logs, types.QueryLogsItem{
			Id:        uint(item.Id),
			Timestamp: item.Timestamp,
			Data:      data,
		})
	}

	return &types.QueryLogsRes{
		Total: rpcRes.Total,
		Logs:  logs,
	}, nil
}
