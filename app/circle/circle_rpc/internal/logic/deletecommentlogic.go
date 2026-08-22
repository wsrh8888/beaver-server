package logic

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(in *circle_rpc.DeleteCommentReq) (*circle_rpc.DeleteCommentRes, error) {
	var comment circle_models.CircleCommentModel
	if err := l.svcCtx.DB.Where("comment_id = ? AND is_deleted = false", in.CommentId).First(&comment).Error; err != nil {
		return nil, fmt.Errorf("评论不存在")
	}

	if err := l.svcCtx.DB.Model(&comment).Update("is_deleted", true).Error; err != nil {
		return nil, fmt.Errorf("删除评论失败: %v", err)
	}

	return &circle_rpc.DeleteCommentRes{}, nil
}
