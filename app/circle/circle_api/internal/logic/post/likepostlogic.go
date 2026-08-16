package post

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikePostLogic {
	return &LikePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikePostLogic) LikePost(req *types.LikePostReq) (resp *types.LikePostRes, err error) {
	var p circle_models.CirclePostModel
	if err = l.svcCtx.DB.Where("post_id = ? AND is_deleted = false", req.PostID).First(&p).Error; err != nil {
		return nil, fmt.Errorf("帖子不存在")
	}

	var like circle_models.CircleLikeModel
	exists := l.svcCtx.DB.Where("post_id = ? AND user_id = ?", req.PostID, req.UserID).First(&like).Error == nil

	if req.Status && !exists {
		newLike := circle_models.CircleLikeModel{
			PostID:   req.PostID,
			UserID:   req.UserID,
			CircleID: p.CircleID,
		}
		if err = l.svcCtx.DB.Create(&newLike).Error; err != nil {
			return nil, fmt.Errorf("点赞失败: %v", err)
		}
	} else if !req.Status && exists {
		if err = l.svcCtx.DB.Delete(&like).Error; err != nil {
			return nil, fmt.Errorf("取消点赞失败: %v", err)
		}
	}

	return &types.LikePostRes{}, nil
}
