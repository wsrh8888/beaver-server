package logic

import (
	"context"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCommentsLogic) ListComments(in *circle_rpc.ListCommentsReq) (*circle_rpc.ListCommentsRes, error) {
	query := l.svcCtx.DB.Model(&circle_models.CircleCommentModel{}).
		Where("post_id = ? AND is_deleted = false", in.PostId)

	var total int64
	query.Count(&total)

	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	var comments []circle_models.CircleCommentModel
	l.svcCtx.DB.Where("post_id = ? AND is_deleted = false", in.PostId).
		Order("created_at ASC").
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&comments)

	list := make([]*circle_rpc.CommentItem, 0, len(comments))
	for _, c := range comments {
		list = append(list, &circle_rpc.CommentItem{
			CommentId: c.CommentID,
			PostId:    c.PostID,
			UserId:    c.UserID,
			Content:   c.Content,
			IsDeleted: c.IsDeleted,
			CreatedAt: c.CreatedAt.String(),
		})
	}

	return &circle_rpc.ListCommentsRes{Total: total, List: list}, nil
}
