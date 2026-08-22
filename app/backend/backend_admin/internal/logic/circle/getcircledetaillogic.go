package circle

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCircleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleDetailLogic {
	return &GetCircleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCircleDetailLogic) GetCircleDetail(req *types.GetCircleDetailReq) (resp *types.GetCircleDetailRes, err error) {
	if req.CircleId == "" {
		return nil, errors.New("圈子ID不能为空")
	}

	rpcRes, err := l.svcCtx.CircleRpc.GetCircleList(l.ctx, &circle_rpc.GetCircleListReq{
		CircleId: req.CircleId,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		l.Errorf("获取圈子详情失败: %v", err)
		return nil, err
	}
	if len(rpcRes.List) == 0 {
		return nil, errors.New("圈子不存在")
	}

	c := mapCircleItem(rpcRes.List[0])
	return &types.GetCircleDetailRes{
		CircleId:    c.CircleId,
		Name:        c.Name,
		Description: c.Description,
		Avatar:      c.Avatar,
		CreatorId:   c.CreatorId,
		JoinType:    c.JoinType,
		MemberCount: c.MemberCount,
		PostCount:   c.PostCount,
		IsDeleted:   c.IsDeleted,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}, nil
}
