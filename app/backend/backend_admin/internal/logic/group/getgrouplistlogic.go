package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupListLogic {
	return &GetGroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupListLogic) GetGroupList(req *types.GetGroupListReq) (resp *types.GetGroupListRes, err error) {
	rpcRes, err := l.svcCtx.GroupRpc.ListGroups(l.ctx, &group_rpc.ListGroupsReq{
		Page:     int32(req.Page),
		PageSize: int32(req.Limit),
		Status:   int32(req.Status),
		Type:     int32(req.Type),
		Keywords: req.Keywords,
	})
	if err != nil {
		l.Errorf("获取群组列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetGroupListItem, 0, len(rpcRes.List))
	for _, g := range rpcRes.List {
		list = append(list, types.GetGroupListItem{
			Id:        uint(g.Id),
			GroupId:   g.GroupId,
			Type:      int(g.Type),
			Title:     g.Title,
			FileName:  g.Avatar,
			CreatorId: g.CreatorId,
			Notice:    g.Notice,
			Status:    int(g.Status),
			CreatedAt: g.CreatedAt,
			UpdatedAt: g.UpdatedAt,
		})
	}
	return &types.GetGroupListRes{List: list, Total: rpcRes.Total}, nil
}
