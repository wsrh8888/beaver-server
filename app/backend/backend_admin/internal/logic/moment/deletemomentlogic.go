package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteMomentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除动态
func NewDeleteMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMomentLogic {
	return &DeleteMomentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMomentLogic) DeleteMoment(req *types.DeleteMomentReq) (resp *types.DeleteMomentRes, err error) {
	// 检查动态是否存在
	var moment moment_models.MomentModel
	err = l.svcCtx.DB.Where("moment_id = ?", req.MomentId).First(&moment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("动态不存在: %s", req.MomentId)
			return nil, errors.New("动态不存在")
		}
		logx.Errorf("查询动态失败: %v", err)
		return nil, errors.New("查询动态失败")
	}

	// 逻辑删除动态（设置is_deleted为true）
	err = l.svcCtx.DB.Model(&moment).Update("is_deleted", true).Error
	if err != nil {
		logx.Errorf("删除动态失败: %v", err)
		return nil, errors.New("删除动态失败")
	}

	// 可选：同时删除相关的评论和点赞
	// 删除评论
	err = l.svcCtx.DB.Where("moment_id = ?", req.MomentId).Delete(&moment_models.MomentCommentModel{}).Error
	if err != nil {
		logx.Errorf("删除动态评论失败: %v", err)
	}

	// 删除点赞
	err = l.svcCtx.DB.Where("moment_id = ?", req.MomentId).Delete(&moment_models.MomentLikeModel{}).Error
	if err != nil {
		logx.Errorf("删除动态点赞失败: %v", err)
	}

	return &types.DeleteMomentRes{}, nil
}
