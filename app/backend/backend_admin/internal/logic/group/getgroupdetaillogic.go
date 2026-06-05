package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupDetailLogic {
	return &GetGroupDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupDetailLogic) GetGroupDetail(req *types.GetGroupDetailReq) (resp *types.GetGroupDetailRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("群组ID不能为空")
	}

	rpcRes, err := l.svcCtx.GroupRpc.ListGroups(l.ctx, &group_rpc.ListGroupsReq{
		Id:       uint64(req.Id),
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		l.Errorf("获取群组详情失败: %v", err)
		return nil, err
	}
	if len(rpcRes.List) == 0 {
		return nil, errors.New("群组不存在")
	}
	g := rpcRes.List[0]
	return &types.GetGroupDetailRes{
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
	}, nil
}
