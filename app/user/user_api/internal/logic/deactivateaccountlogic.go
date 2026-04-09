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
)

type DeactivateAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注销账号
func NewDeactivateAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeactivateAccountLogic {
	return &DeactivateAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeactivateAccountLogic) DeactivateAccount(req *types.DeactivateAccountReq) (resp *types.DeactivateAccountRes, err error) {
	// 查询用户
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ? AND status = 1", req.UserID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码，防止误操作
	if !pwd.CheckPad(user.Password, req.Password) {
		return nil, errors.New("密码错误")
	}

	// 标记账号为已注销（status=3）
	if err := l.svcCtx.DB.Model(&user).Update("status", 3).Error; err != nil {
		return nil, errors.New("注销失败")
	}

	// 删除所有设备类型的登录态，强制全端下线
	deviceTypes := []string{"desktop", "mobile", "web", "unknown"}
	for _, dt := range deviceTypes {
		key := fmt.Sprintf("login_%s_%s", req.UserID, dt)
		l.svcCtx.Redis.Del(key)
	}

	return &types.DeactivateAccountRes{}, nil
}
