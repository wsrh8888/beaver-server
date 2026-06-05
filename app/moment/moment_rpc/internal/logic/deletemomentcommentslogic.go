package logic

import (
	"context"

	"beaver/app/moment/moment_models"
	"beaver/app/moment/moment_rpc/internal/svc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMomentCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMomentCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMomentCommentsLogic {
	return &DeleteMomentCommentsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteMomentCommentsLogic) DeleteMomentComments(in *moment_rpc.DeleteMomentCommentsReq) (*moment_rpc.DeleteMomentCommentsRes, error) {
	if len(in.CommentIds) == 0 {
		return &moment_rpc.DeleteMomentCommentsRes{}, nil
	}
	if err := l.svcCtx.DB.Where("comment_id IN ?", in.CommentIds).Delete(&moment_models.MomentCommentModel{}).Error; err != nil {
		l.Errorf("删除评论失败: %v", err)
		return nil, err
	}
	return &moment_rpc.DeleteMomentCommentsRes{}, nil
}
