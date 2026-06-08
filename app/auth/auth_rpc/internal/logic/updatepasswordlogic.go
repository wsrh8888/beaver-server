package logic

import (
	"context"
	"errors"

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)


type UpdatePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		ctx:    ctx,
		logger: logger.New("update_password"),
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(in *auth_rpc.UpdatePasswordReq) (*auth_rpc.UpdatePasswordRes, error) {
	// 验证必填字段
	if in.UserId == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if in.OldPassword == "" || in.NewPassword == "" {
		return nil, errors.New("旧密码和新密码不能为空")
	}

	// 查询用户凭证
	var credential auth_models.AuthCredentialModel
	err := l.svcCtx.DB.Take(&credential, "user_id = ?", in.UserId).Error
	if err != nil {
		logx.Errorf("查询用户凭证失败: %v", err)
		return nil, errors.New("用户凭证不存在")
	}

	// 验证旧密码
	if !pwd.CheckPad(credential.Password, in.OldPassword) {
		return nil, errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword := pwd.HahPwd(in.NewPassword)

	// 更新密码
	credential.Password = hashedPassword
	err = l.svcCtx.DB.Save(&credential).Error
	if err != nil {
		logx.Errorf("更新密码失败: %v", err)
		return nil, errors.New("更新密码失败")
	}

	logx.Infof("密码更新成功: userID=%s", in.UserId)
	l.logger.Info(model.LogMsg{
		Text: "密码修改成功",
		Data: map[string]interface{}{
			"userId": in.UserId,
		},
	})

	return &auth_rpc.UpdatePasswordRes{
		Success: true,
	}, nil
}
