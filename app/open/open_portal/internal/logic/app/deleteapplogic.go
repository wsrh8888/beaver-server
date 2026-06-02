package app

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除应用
func NewDeleteAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAppLogic {
	return &DeleteAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAppLogic) DeleteApp(req *types.DeleteAppReq) (resp *types.DeleteAppRes, err error) {
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// 软删除应用（GORM 的 Delete 会设置 deleted_at）
	result := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).Delete(&open_models.OpenApp{})
	if result.Error != nil {
		logx.Errorf("删除应用失败: %v", result.Error)
		return nil, errors.New("删除失败")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("应用不存在或无权限")
	}

	logx.Infof("应用删除成功: app_id=%s", req.AppID)

	return &types.DeleteAppRes{}, nil
}
