package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventListLogic {
	return &GetEventListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventListLogic) GetEventList(req *types.GetEventListReq) (resp *types.GetEventListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.AdminGetEventList(l.ctx, &platform_rpc.AdminGetEventListReq{
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
		BucketId:   req.BucketID,
		EventName:  req.EventName,
		Action:     req.Action,
		UserFilter: req.UserFilter,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Platform:   req.Platform,
	})
	if err != nil {
		l.Errorf("获取事件列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetEventListItem, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		list = append(list, types.GetEventListItem{
			Id:         uint(item.Id),
			EventName:  item.EventName,
			Action:     item.Action,
			UserID:     item.UserId,
			BucketID:   item.BucketId,
			BucketName: item.BucketName,
			Platform:   item.Platform,
			DeviceID:   item.DeviceId,
			Data:       item.Data,
			Timestamp:  item.Timestamp,
			CreatedAt:  item.CreatedAt,
		})
	}

	return &types.GetEventListRes{
		List:  list,
		Total: rpcRes.Total,
	}, nil
}
