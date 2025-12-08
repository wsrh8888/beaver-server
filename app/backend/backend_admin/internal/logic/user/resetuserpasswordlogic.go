package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_models"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ResetUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置用户密码
func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetUserPasswordLogic) ResetUserPassword(req *types.ResetUserPasswordReq) (resp *types.ResetUserPasswordRes, err error) {
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

	// 验证新密码不能与原密码相同
	if pwd.CheckPad(user.Password, req.NewPassword) {
		return nil, errors.New("新密码不能与原密码相同")
	}

	// 加密新密码
	hashedPassword := pwd.HahPwd(req.NewPassword)

	// 更新密码
	err = l.svcCtx.DB.Model(&user).Update("password", hashedPassword).Error
	if err != nil {
		logx.Errorf("重置用户密码失败: %v", err)
		return nil, errors.New("重置用户密码失败")
	}

	return &types.ResetUserPasswordRes{}, nil
}
