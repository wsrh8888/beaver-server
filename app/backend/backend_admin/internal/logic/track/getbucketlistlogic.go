package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBucketListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBucketListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBucketListLogic {
	return &GetBucketListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBucketListLogic) GetBucketList(req *types.GetBucketListReq) (resp *types.GetBucketListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.AdminGetBucketList(l.ctx, &platform_rpc.AdminGetBucketListReq{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		Keyword:  req.Keyword,
		IsActive: req.IsActive,
	})
	if err != nil {
		l.Errorf("获取 Bucket 列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetBucketListItem, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		list = append(list, types.GetBucketListItem{
			BucketId:    item.BucketId,
			Name:        item.Name,
			Description: item.Description,
			CreateUser:  item.CreateUser,
			IsActive:    item.IsActive,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	return &types.GetBucketListRes{
		List:  list,
		Total: rpcRes.Total,
	}, nil
}
