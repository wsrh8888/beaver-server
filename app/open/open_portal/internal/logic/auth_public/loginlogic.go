package auth_public

import (
	"context"
	"errors"
	"time"

	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/jwts"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 开发者 Portal 登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRes, err error) {
	userRes, err := l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
		Keyword: req.Username,
		Type:    "phone",
	})
	if err != nil {
		userRes, err = l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
			Keyword: req.Username,
			Type:    "email",
		})
	}
	if err != nil {
		logx.Errorf("用户不存在: %s, error: %v", req.Username, err)
		return nil, errors.New("用户名或密码错误")
	}

	verifyRes, err := l.svcCtx.AuthRpc.VerifyPassword(l.ctx, &auth_rpc.VerifyPasswordReq{
		UserId:   userRes.UserInfo.UserId,
		Password: req.Password,
	})
	if err != nil || !verifyRes.Valid {
		logx.Errorf("密码错误: user_id=%s", userRes.UserInfo.UserId)
		return nil, errors.New("用户名或密码错误")
	}

	secretKey := l.svcCtx.Config.Auth.AccessSecret
	expireHours := l.svcCtx.Config.Auth.AccessExpire / 3600
	if expireHours == 0 {
		expireHours = 12
	}

	token, err := jwts.GenToken(jwts.JwtPayLoad{
		UserID:   userRes.UserInfo.UserId,
		NickName: userRes.UserInfo.NickName,
	}, secretKey, int(expireHours))
	if err != nil {
		logx.Errorf("生成 token 失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	expireAt := time.Now().Add(time.Duration(expireHours) * time.Hour).UnixMilli()

	logx.Infof("开放平台登录成功: user_id=%s, nick_name=%s", userRes.UserInfo.UserId, userRes.UserInfo.NickName)

	return &types.LoginRes{
		Token:    token,
		UserID:   userRes.UserInfo.UserId,
		NickName: userRes.UserInfo.NickName,
		ExpireAt: expireAt,
	}, nil
}
