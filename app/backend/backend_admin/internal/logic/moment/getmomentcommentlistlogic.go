package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_models"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取动态评论列表
func NewGetMomentCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentCommentListLogic {
	return &GetMomentCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMomentCommentListLogic) GetMomentCommentList(req *types.GetMomentCommentListReq) (resp *types.GetMomentCommentListRes, err error) {
	// 分页参数校验
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// 构建查询条件（使用 moment id）
	whereClause := l.svcCtx.DB.Where("moment_id = ?", req.MomentId)

	// 分页查询评论
	comments, count, err := list_query.ListQuery(l.svcCtx.DB, moment_models.MomentCommentModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  page,
			Limit: limit,
			Sort:  "created_at desc",
		},
		Where: whereClause,
	})

	if err != nil {
		logx.Errorf("查询动态评论列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.GetMomentCommentListItem
	for _, comment := range comments {
		list = append(list, types.GetMomentCommentListItem{
			CommentId: comment.CommentID,
			MomentId:  comment.MomentID,
			UserId:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.String(),
			UpdatedAt: comment.UpdatedAt.String(),
		})
	}

	return &types.GetMomentCommentListRes{
		List:  list,
		Total: count,
	}, nil
}
