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
	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("moment_id = ?", req.MomentId)

	// 分页查询评论
	comments, count, err := list_query.ListQuery(l.svcCtx.DB, moment_models.MomentCommentModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where: whereClause,
	})

	if err != nil {
		logx.Errorf("查询动态评论列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.MomentCommentInfo
	for _, comment := range comments {
		list = append(list, types.MomentCommentInfo{
			Id:        comment.Id,
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
