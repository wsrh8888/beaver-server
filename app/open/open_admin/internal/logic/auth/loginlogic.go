package auth

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_admin/internal/middleware"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"

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
	// TODO: 实际应该查询用户表验证用户名密码
	// 这里简化处理，假设 admin/admin123 是管理员
	if req.Username != "admin" || req.Password != "admin123" {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 JWT Token
	secretKey := l.svcCtx.Config.JWT.SecretKey
	expireHours := l.svcCtx.Config.JWT.ExpireHours
	if expireHours == 0 {
		expireHours = 24 // 默认 24 小时
	}

	token, err := middleware.GenerateJWT("admin", "admin", secretKey, expireHours)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	expireAt := time.Now().Add(time.Duration(expireHours) * time.Hour).UnixMilli()

	return &types.LoginRes{
		Token:    token,
		UserID:   "admin",
		Role:     "admin",
		ExpireAt: expireAt,
	}, nil
}
