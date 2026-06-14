package logic

import (
	"context"
	"fmt"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)


type DeleteMomentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewDeleteMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMomentLogic {
	return &DeleteMomentLogic{
		ctx:    ctx,
		logger: logger.New("delete_moment"),
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
	l.logger.Info(model.LogMsg{
		Text: "朋友圈删除成功",
		Data: map[string]interface{}{
			"momentId": req.MomentID,
			"userId":   req.UserID,
		},
	})
	return resp, nil
}
