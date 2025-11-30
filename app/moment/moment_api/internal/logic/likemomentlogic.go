package logic

import (
	"context"
	"errors"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LikeMomentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeMomentLogic {
	return &LikeMomentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeMomentLogic) LikeMoment(req *types.LikeMomentReq) (resp *types.LikeMomentRes, err error) {
	var like moment_models.MomentLikeModel

	if req.Status {
		// 点赞操作
		like = moment_models.MomentLikeModel{
			UUID:     uuid.New().String(),
			MomentID: req.MomentID,
			UserID:   req.UserID,
		}

		// 检查是否已点赞
		result := l.svcCtx.DB.Where("moment_id = ? AND user_id = ?", req.MomentID, req.UserID).First(&like)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 未找到记录，进行点赞
			if err := l.svcCtx.DB.Create(&like).Error; err != nil {
				return nil, err
			}
		} else if result.Error != nil {
			// 其他错误
			return nil, result.Error
		} else {
			// 记录存在，如果被软删除了则恢复
			if like.IsDeleted {
				if err := l.svcCtx.DB.Model(&like).Update("is_deleted", false).Error; err != nil {
					return nil, err
				}
			}
			// 如果记录存在且未被删除，则无需操作
		}
	} else {
		// 取消点赞操作（软删除）
		if err := l.svcCtx.DB.Model(&moment_models.MomentLikeModel{}).Where("moment_id = ? AND user_id = ? AND is_deleted = false", req.MomentID, req.UserID).Update("is_deleted", true).Error; err != nil {
			return nil, err
		}
	}

	return resp, nil
}
