package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
	// 检查用户是否存在
	var user user_models.UserModel
	err = l.svcCtx.DB.Where("user_id = ?", req.UserID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("用户不存在: %s", req.UserID)
			return nil, errors.New("用户不存在")
		}
		logx.Errorf("查询用户失败: %v", err)
		return nil, errors.New("查询用户失败")
	}

	// 逻辑删除用户（设置状态为删除状态）
	err = l.svcCtx.DB.Model(&user).Update("status", 3).Error
	if err != nil {
		logx.Errorf("删除用户失败: %v", err)
		return nil, errors.New("删除用户失败")
	}

	// 或者使用物理删除（根据业务需求选择）
	// err = l.svcCtx.DB.Delete(&user).Error
	// if err != nil {
	//     logx.Errorf("删除用户失败: %v", err)
	//     return nil, errors.New("删除用户失败")
	// }

	return &types.DeleteUserRes{}, nil
}
