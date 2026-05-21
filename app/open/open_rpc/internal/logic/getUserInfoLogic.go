package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *open_rpc.GetUserInfoReq) (*open_rpc.GetUserInfoRes, error) {
	// 1. 查询 Access Token
	var token open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("token = ?", in.AccessToken).First(&token).Error; err != nil {
		return nil, errors.New("无效的访问令牌")
	}

	// 2. 检查 Token 是否过期
	if time.Now().Unix() > token.ExpiresAt {
		return nil, errors.New("访问令牌已过期")
	}

	// 3. 查询用户信息
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ?", token.UserID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return &open_rpc.GetUserInfoRes{
		UserId:   user.UserID,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Phone:    user.Phone,
		Email:    user.Email,
	}, nil
}
