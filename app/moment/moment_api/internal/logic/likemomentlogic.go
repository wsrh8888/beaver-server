package logic

import (
	"context"
	"errors"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"

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
		}
	} else {
		// 取消点赞操作
		if err := l.svcCtx.DB.Where("moment_id = ? AND user_id = ?", req.MomentID, req.UserID).Delete(&moment_models.MomentLikeModel{}).Error; err != nil {
			return nil, err
		}
	}

	return resp, nil
}
