package circle

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCircleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleListLogic {
	return &GetCircleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCircleListLogic) GetCircleList(req *types.GetCircleListReq) (resp *types.GetCircleListRes, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	rpcRes, err := l.svcCtx.CircleRpc.GetCircleList(l.ctx, &circle_rpc.GetCircleListReq{
		Page:     int32(page),
		PageSize: int32(limit),
		UserId:   req.UserId,
		Keywords: req.Keywords,
		CircleId: req.CircleId,
	})
	if err != nil {
		l.Errorf("获取圈子列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetCircleListItem, 0, len(rpcRes.List))
	for _, c := range rpcRes.List {
		list = append(list, mapCircleItem(c))
	}
	return &types.GetCircleListRes{List: list, Total: rpcRes.Total}, nil
}

func mapCircleItem(c *circle_rpc.CircleItem) types.GetCircleListItem {
	if c == nil {
		return types.GetCircleListItem{}
	}
	return types.GetCircleListItem{
		CircleId:    c.CircleId,
		Name:        c.Name,
		Description: c.Description,
		Avatar:      c.Avatar,
		CreatorId:   c.CreatorId,
		JoinType:    int(c.JoinType),
		MemberCount: c.MemberCount,
		PostCount:   c.PostCount,
		IsDeleted:   c.IsDeleted,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
