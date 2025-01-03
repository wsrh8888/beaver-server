package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdatePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(req *types.UpdatePasswordReq) (resp *types.UpdatePasswordRes, err error) {
	// 查询用户信息
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ?", req.UserId).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	if !pwd.CheckPad(user.Password, req.OldPassword) {
		fmt.Println("原始密码错误", user.Password, req.OldPassword)
		return nil, errors.New("原始密码错误")
	}
	if pwd.CheckPad(user.Password, req.NewPassword) {
		return nil, errors.New("不能与原密码相同")
	}
	if err := l.svcCtx.DB.Model(&user).Update("password", pwd.HahPwd(req.NewPassword)).Error; err != nil {
		return nil, err
	}
	return nil, nil
}
