package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/utils"
	"beaver/utils/jwts"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthenticationLogic) Authentication(req *types.AuthenticationReq) (resp *types.AuthenticationRes, err error) {
	if utils.InListByRegex(l.svcCtx.Config.WhiteList, req.ValidPath) {
		logx.Infof("白名单请求：%s, %s", req.ValidPath, req.Token)
		return
	}
	if req.Token == "" {
		err = errors.New("token不能为空")
		return
	}
	claims, err := jwts.ParseToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		err = errors.New("认证失败")
		return
	}
	key := fmt.Sprintf("login_%s", claims.UserID)
	token, _ := l.svcCtx.Redis.Get(key).Result()
	if token != req.Token {
		fmt.Println("token不一致", token, req.Token)
		err = errors.New("token已失效")
		return
	}
	resp = &types.AuthenticationRes{
		UserID: claims.UserID,
	}
	fmt.Println(resp, "数据")
	return resp, nil
}
