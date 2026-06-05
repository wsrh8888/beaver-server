package logic

import (
	"context"
	"errors"

	"beaver/app/moment/moment_models"
	"beaver/app/moment/moment_rpc/internal/svc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateMomentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMomentLogic {
	return &UpdateMomentLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateMomentLogic) UpdateMoment(in *moment_rpc.UpdateMomentReq) (*moment_rpc.UpdateMomentRes, error) {
	var moment moment_models.MomentModel
	if err := l.svcCtx.DB.Where("moment_id = ?", in.MomentId).First(&moment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("动态不存在")
		}
		return nil, err
	}

	if in.IsDeleted != nil && *in.IsDeleted {
		if err := l.svcCtx.DB.Model(&moment).Update("is_deleted", true).Error; err != nil {
			return nil, err
		}
		_ = l.svcCtx.DB.Where("moment_id = ?", in.MomentId).Delete(&moment_models.MomentCommentModel{}).Error
		_ = l.svcCtx.DB.Where("moment_id = ?", in.MomentId).Delete(&moment_models.MomentLikeModel{}).Error
	}
	return &moment_rpc.UpdateMomentRes{}, nil
}
