package contact

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	user_models "beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除用户
func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.DeleteUserRes, err error) {
	// 1. 查询用户是否存在
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 软删除用户（更新状态为删除）
	if err := l.svcCtx.DB.Model(&user).Update("status", 3).Error; err != nil {
		logx.Errorf("删除用户失败: %v", err)
		return nil, errors.New("删除用户失败")
	}

	return &types.DeleteUserRes{
		Success: true,
	}, nil
}
