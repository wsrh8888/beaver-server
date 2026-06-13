package auth

import (
	"context"
	"errors"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
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

func (l *UpdatePasswordLogic) UpdatePassword(req *types.UpdatePasswordReq) (*types.UpdatePasswordRes, error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if req.OldPassword == "" || req.NewPassword == "" {
		return nil, errors.New("旧密码和新密码不能为空")
	}
	if len(req.NewPassword) < 6 {
		return nil, errors.New("新密码长度不能少于6位")
	}

	var credential auth_models.AuthCredentialModel
	if err := l.svcCtx.DB.Take(&credential, "user_id = ?", req.UserID).Error; err != nil {
		logx.Errorf("查询用户凭证失败: %v", err)
		return nil, errors.New("用户凭证不存在")
	}

	if !pwd.CheckPad(credential.Password, req.OldPassword) {
		return nil, errors.New("旧密码错误")
	}

	credential.Password = pwd.HahPwd(req.NewPassword)
	if err := l.svcCtx.DB.Save(&credential).Error; err != nil {
		logx.Errorf("更新密码失败: %v", err)
		return nil, errors.New("更新密码失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "密码修改成功",
		Data: map[string]interface{}{"userId": req.UserID},
	})
	return &types.UpdatePasswordRes{}, nil
}
