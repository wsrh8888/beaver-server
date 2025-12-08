package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除用户
func NewBatchDeleteUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteUsersLogic {
	return &BatchDeleteUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteUsersLogic) BatchDeleteUsers(req *types.BatchDeleteUsersReq) (resp *types.BatchDeleteUsersRes, err error) {
	// 批量逻辑删除用户（设置状态为删除状态）
	err = l.svcCtx.DB.Model(&user_models.UserModel{}).
		Where("user_id IN ?", req.Ids).
		Update("status", 3).Error

	if err != nil {
		logx.Errorf("批量删除用户失败: %v", err)
		return nil, errors.New("批量删除用户失败")
	}

	// 或者使用物理删除（根据业务需求选择）
	// err = l.svcCtx.DB.Where("user_id IN ?", req.Ids).Delete(&user_models.UserModel{}).Error
	// if err != nil {
	//     logx.Errorf("批量删除用户失败: %v", err)
	//     return nil, errors.New("批量删除用户失败")
	// }

	return &types.BatchDeleteUsersRes{}, nil
}
