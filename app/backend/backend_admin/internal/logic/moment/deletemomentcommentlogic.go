package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMomentCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMomentCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMomentCommentLogic {
	return &DeleteMomentCommentLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteMomentCommentLogic) DeleteMomentComment(req *types.DeleteMomentCommentReq) (resp *types.DeleteMomentCommentRes, err error) {
	if req.CommentId == "" {
		return nil, errors.New("评论ID不能为空")
	}

	_, err = l.svcCtx.MomentRpc.DeleteMomentComments(l.ctx, &moment_rpc.DeleteMomentCommentsReq{
		CommentIds: []string{req.CommentId},
	})
	if err != nil {
		l.Errorf("删除评论失败: %v", err)
		return nil, err
	}
	return &types.DeleteMomentCommentRes{}, nil
}
