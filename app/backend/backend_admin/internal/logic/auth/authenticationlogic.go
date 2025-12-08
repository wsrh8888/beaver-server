package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	"beaver/utils"
	"beaver/utils/jwts"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户认证
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
		return &types.AuthenticationRes{}, nil
	}

	claims, err := jwts.ParseToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		err = errors.New("认证失败")
		return
	}

	// 验证管理员用户状态
	var adminUser backend_models.AdminUser
	err = l.svcCtx.DB.Take(&adminUser, "user_id = ? AND status = ?", claims.UserID, 1).Error
	if err != nil {
		err = errors.New("管理员用户不存在或已被禁用")
		return
	}

	key := fmt.Sprintf("admin_login_%s", claims.UserID)
	token, _ := l.svcCtx.Redis.Get(key).Result()
	if token != req.Token {
		fmt.Println("token不一致", token, req.Token)
		err = errors.New("token已失效， token不一致")
		return
	}
	resp = &types.AuthenticationRes{
		UserID: claims.UserID,
	}
	fmt.Println(resp, "数据")
	return resp, nil
}
