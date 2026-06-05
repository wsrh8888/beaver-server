package logic

import (
	"context"
	"time"

	"beaver/app/moment/moment_models"
	"beaver/app/moment/moment_rpc/internal/svc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMomentCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMomentCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMomentCommentsLogic {
	return &ListMomentCommentsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListMomentCommentsLogic) ListMomentComments(in *moment_rpc.ListMomentCommentsReq) (*moment_rpc.ListMomentCommentsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&moment_models.MomentCommentModel{})
	if in.MomentId != "" {
		db = db.Where("moment_id = ?", in.MomentId)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计评论失败: %v", err)
		return nil, err
	}

	var list []moment_models.MomentCommentModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询评论列表失败: %v", err)
		return nil, err
	}

	items := make([]*moment_rpc.MomentCommentItem, 0, len(list))
	for _, c := range list {
		items = append(items, &moment_rpc.MomentCommentItem{
			CommentId: c.CommentID,
			MomentId:  c.MomentID,
			UserId:    c.UserID,
			Content:   c.Content,
			CreatedAt: time.Time(c.CreatedAt).Format(time.RFC3339),
			UpdatedAt: time.Time(c.UpdatedAt).Format(time.RFC3339),
		})
	}
	return &moment_rpc.ListMomentCommentsRes{Total: total, List: items}, nil
}
