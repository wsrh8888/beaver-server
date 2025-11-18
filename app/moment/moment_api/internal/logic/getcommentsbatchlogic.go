package logic

import (
	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentsBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取评论数据（用于数据同步）
func NewGetCommentsBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentsBatchLogic {
	return &GetCommentsBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentsBatchLogic) GetCommentsBatch(req *types.GetCommentsBatchReq) (resp *types.GetCommentsBatchRes, err error) {
	if len(req.UserIds) == 0 {
		return &types.GetCommentsBatchRes{
			Comments:    []types.CommentBatchItem{},
			HasMore:     false,
			NextVersion: req.EndVersion,
		}, nil
	}

	var comments []moment_models.MomentCommentModel
	query := l.svcCtx.DB.Where("moment_user_id IN (?)", req.UserIds).
		Where("version > ? AND version <= ?", req.StartVersion, req.EndVersion).
		Order("version ASC").
		Limit(req.Limit + 1) // 多查一条来判断是否有更多数据

	err = query.Find(&comments).Error
	if err != nil {
		l.Errorf("查询评论数据失败: %v", err)
		return nil, err
	}

	// 检查是否有更多数据
	hasMore := len(comments) > req.Limit
	if hasMore {
		comments = comments[:req.Limit] // 移除多查的那条
	}

	// 转换为响应格式
	var commentItems []types.CommentBatchItem
	for _, comment := range comments {
		commentItems = append(commentItems, types.CommentBatchItem{
			UUID:         comment.UUID,
			MomentID:     comment.MomentID,
			UserID:       comment.UserID,
			MomentUserID: comment.MomentUserID,
			Content:      comment.Content,
			Version:      comment.Version,
			CreateAt:     time.Time(comment.CreatedAt).UnixMilli(),
			UpdateAt:     time.Time(comment.UpdatedAt).UnixMilli(),
			IsDeleted:    comment.IsDeleted,
		})
	}

	// 计算下次同步的起始版本号
	nextVersion := req.EndVersion
	if len(comments) > 0 {
		nextVersion = comments[len(comments)-1].Version
	}

	return &types.GetCommentsBatchRes{
		Comments:    commentItems,
		HasMore:     hasMore,
		NextVersion: nextVersion,
	}, nil
}
