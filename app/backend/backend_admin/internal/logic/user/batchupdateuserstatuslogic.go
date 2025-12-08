package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量更新用户状态
func NewBatchUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateUserStatusLogic {
	return &BatchUpdateUserStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateUserStatusLogic) BatchUpdateUserStatus(req *types.BatchUpdateUserStatusReq) (resp *types.BatchUpdateUserStatusRes, err error) {
	if len(req.Ids) == 0 {
		return nil, errors.New("请选择要更新的用户")
	}

	// 验证状态值的有效性
	if req.Status < 1 || req.Status > 3 {
		return nil, errors.New("无效的状态值，状态值应为：1-正常，2-禁用，3-删除")
	}

	// 批量更新用户状态
	err = l.svcCtx.DB.Model(&user_models.UserModel{}).
		Where("user_id IN ?", req.Ids).
		Update("status", int8(req.Status)).Error

	if err != nil {
		logx.Errorf("批量更新用户状态失败: %v", err)
		return nil, errors.New("批量更新用户状态失败")
	}

	return &types.BatchUpdateUserStatusRes{}, nil
}
