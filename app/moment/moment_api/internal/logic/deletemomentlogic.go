package logic

import (
	"context"
	"fmt"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMomentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMomentLogic {
	return &DeleteMomentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMomentLogic) DeleteMoment(req *types.DeleteMomentReq) (resp *types.DeleteMomentRes, err error) {
	// 检查动态是否存在以及用户是否有权限删除该动态
	var moment moment_models.MomentModel
	if err := l.svcCtx.DB.Where("id = ? AND user_id = ?", req.MomentID, req.UserID).First(&moment).Error; err != nil {
		return nil, fmt.Errorf("删除失败")
	}

	// 将 is_deleted 标记为 true
	if err := l.svcCtx.DB.Model(&moment).Update("is_deleted", true).Error; err != nil {
		return nil, fmt.Errorf("failed to delete moment: %v", err)
	}

	resp = &types.DeleteMomentRes{}
	return resp, nil
}
