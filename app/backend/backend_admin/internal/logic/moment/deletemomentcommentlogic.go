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

type DeleteMomentCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除动态评论
func NewDeleteMomentCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMomentCommentLogic {
	return &DeleteMomentCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMomentCommentLogic) DeleteMomentComment(req *types.DeleteMomentCommentReq) (resp *types.DeleteMomentCommentRes, err error) {
	// 检查评论是否存在
	var comment moment_models.MomentCommentModel
	err = l.svcCtx.DB.Where("uuid = ?", req.Uuid).First(&comment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("评论不存在: %s", req.Uuid)
			return nil, errors.New("评论不存在")
		}
		logx.Errorf("查询评论失败: %v", err)
		return nil, errors.New("查询评论失败")
	}

	// 删除评论
	err = l.svcCtx.DB.Delete(&comment).Error
	if err != nil {
		logx.Errorf("删除评论失败: %v", err)
		return nil, errors.New("删除评论失败")
	}

	return &types.DeleteMomentCommentRes{}, nil
}
