package circle

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCircleCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCircleCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCircleCommentLogic {
	return &DeleteCircleCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCircleCommentLogic) DeleteCircleComment(req *types.DeleteCircleCommentReq) (resp *types.DeleteCircleCommentRes, err error) {
	if req.CommentId == "" {
		return nil, errors.New("评论ID不能为空")
	}

	_, err = l.svcCtx.CircleRpc.DeleteComment(l.ctx, &circle_rpc.DeleteCommentReq{
		CommentId: req.CommentId,
	})
	if err != nil {
		l.Errorf("删除帖子评论失败: %v", err)
		return nil, err
	}
	return &types.DeleteCircleCommentRes{}, nil
}
