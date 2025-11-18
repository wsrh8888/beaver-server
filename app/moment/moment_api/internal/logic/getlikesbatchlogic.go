package logic

import (
	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLikesBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取点赞数据（用于数据同步）
func NewGetLikesBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLikesBatchLogic {
	return &GetLikesBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLikesBatchLogic) GetLikesBatch(req *types.GetLikesBatchReq) (resp *types.GetLikesBatchRes, err error) {
	if len(req.UserIds) == 0 {
		return &types.GetLikesBatchRes{
			Likes:       []types.LikeBatchItem{},
			HasMore:     false,
			NextVersion: req.EndVersion,
		}, nil
	}

	var likes []moment_models.MomentLikeModel
	query := l.svcCtx.DB.Where("moment_user_id IN (?)", req.UserIds).
		Where("version > ? AND version <= ?", req.StartVersion, req.EndVersion).
		Order("version ASC").
		Limit(req.Limit + 1) // 多查一条来判断是否有更多数据

	err = query.Find(&likes).Error
	if err != nil {
		l.Errorf("查询点赞数据失败: %v", err)
		return nil, err
	}

	// 检查是否有更多数据
	hasMore := len(likes) > req.Limit
	if hasMore {
		likes = likes[:req.Limit] // 移除多查的那条
	}

	// 转换为响应格式
	var likeItems []types.LikeBatchItem
	for _, like := range likes {
		likeItems = append(likeItems, types.LikeBatchItem{
			UUID:         like.UUID,
			MomentID:     like.MomentID,
			UserID:       like.UserID,
			MomentUserID: like.MomentUserID,
			Version:      like.Version,
			CreateAt:     time.Time(like.CreatedAt).UnixMilli(),
			UpdateAt:     time.Time(like.UpdatedAt).UnixMilli(),
			IsDeleted:    like.IsDeleted,
		})
	}

	// 计算下次同步的起始版本号
	nextVersion := req.EndVersion
	if len(likes) > 0 {
		nextVersion = likes[len(likes)-1].Version
	}

	return &types.GetLikesBatchRes{
		Likes:       likeItems,
		HasMore:     hasMore,
		NextVersion: nextVersion,
	}, nil
}
