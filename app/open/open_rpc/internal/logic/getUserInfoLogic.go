package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

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
	var token open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("token = ?", in.AccessToken).First(&token).Error; err != nil {
		return nil, errors.New("无效的访问令牌")
	}
	if time.Now().Unix() > token.ExpiresAt {
		return nil, errors.New("访问令牌已过期")
	}

	userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: token.UserID})
	if err != nil || userRes.UserInfo == nil {
		return nil, errors.New("用户不存在")
	}
	info := userRes.UserInfo

	return &open_rpc.GetUserInfoRes{
		UserId:   info.UserId,
		NickName: info.NickName,
		Avatar:   info.Avatar,
		Phone:    info.Phone,
		Email:    info.Email,
	}, nil
}
