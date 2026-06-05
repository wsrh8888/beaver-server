package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMomentCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentCommentListLogic {
	return &GetMomentCommentListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetMomentCommentListLogic) GetMomentCommentList(req *types.GetMomentCommentListReq) (resp *types.GetMomentCommentListRes, err error) {
	if req.MomentId == "" {
		return nil, errors.New("动态ID不能为空")
	}

	rpcRes, err := l.svcCtx.MomentRpc.ListMomentComments(l.ctx, &moment_rpc.ListMomentCommentsReq{
		MomentId: req.MomentId,
		Page:     int32(req.Page),
		PageSize: int32(req.Limit),
	})
	if err != nil {
		l.Errorf("获取动态评论列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetMomentCommentListItem, 0, len(rpcRes.List))
	for _, c := range rpcRes.List {
		list = append(list, types.GetMomentCommentListItem{
			CommentId: c.CommentId,
			MomentId:  c.MomentId,
			UserId:    c.UserId,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}
	return &types.GetMomentCommentListRes{List: list, Total: rpcRes.Total}, nil
}
