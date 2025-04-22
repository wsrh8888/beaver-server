package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/utils/jwts"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRes, err error) {
	// 查询用户信息
	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "phone = ?", req.Phone).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if !pwd.CheckPad(user.Password, req.Password) {
		return nil, errors.New("密码错误")
	}

	// 生成token
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		Phone:    user.Phone,
		Nickname: user.NickName,
		UserID:   user.UUID,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Errorf("生成token失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 存储token到Redis
	key := fmt.Sprintf("login_%s", user.UUID)
	err = l.svcCtx.Redis.Set(key, token, time.Hour*48).Err()
	if err != nil {
		logx.Errorf("存储token失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 记录登录设备信息
	deviceKey := fmt.Sprintf("device_%s", user.UUID)
	deviceInfo := map[string]interface{}{
		"user_agent": l.ctx.Value("user-agent"),
		"ip":         l.ctx.Value("client-ip"),
		"login_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := l.svcCtx.Redis.HMSet(deviceKey, deviceInfo).Err(); err != nil {
		logx.Errorf("记录设备信息失败: %v", err)
	}

	return &types.LoginRes{
		Token:  token,
		UserID: user.UUID,
	}, nil
}
