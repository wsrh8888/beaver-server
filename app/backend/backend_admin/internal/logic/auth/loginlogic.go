package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	"beaver/utils/jwts"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRes, err error) {
	// 直接查询管理员用户
	var adminUser backend_models.AdminUser

	// 根据手机号查找管理员用户
	err = l.svcCtx.DB.Take(&adminUser, "phone = ? AND status = ?", req.Phone, 1).Error
	if err != nil {
		// 如果是超级管理员手机号且不存在，则创建
		if req.Phone == "15383645663" {
			adminUser = backend_models.AdminUser{
				UUID:     fmt.Sprintf("admin_%d", time.Now().Unix()),
				NickName: "超级管理员",
				Password: pwd.HahPwd(req.Password),
				Phone:    req.Phone,
				Status:   1,
			}
			err = l.svcCtx.DB.Create(&adminUser).Error
			if err != nil {
				logx.Errorf("创建超级管理员失败: %v", err)
				return nil, errors.New("服务内部异常")
			}
			return nil, errors.New("管理员账户创建成功，请重新登录")
		}
		return nil, errors.New("管理员用户不存在")
	}

	// 验证密码
	if !pwd.CheckPad(adminUser.Password, req.Password) {
		return nil, errors.New("密码错误")
	}

	// 检查用户状态
	if adminUser.Status != 1 {
		return nil, errors.New("账户已被禁用")
	}

	// 生成JWT token
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		Nickname: adminUser.NickName,
		UserID:   adminUser.UUID,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Errorf("生成 token 失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 存储token到Redis
	key := fmt.Sprintf("admin_login_%s", adminUser.UUID)
	err = l.svcCtx.Redis.Set(key, token, time.Hour*48).Err()
	if err != nil {
		logx.Errorf("存储 token 失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 返回登录结果
	return &types.LoginRes{
		Token:  token,
		UserID: adminUser.UUID,
	}, nil
}
