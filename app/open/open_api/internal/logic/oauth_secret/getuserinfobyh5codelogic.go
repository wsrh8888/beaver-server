package oauth_secret

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/logic/oauthutil"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoByH5CodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoByH5CodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByH5CodeLogic {
	return &GetUserInfoByH5CodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoByH5CodeLogic) GetUserInfoByH5Code(req *types.GetUserInfoByH5CodeReq) (resp *types.GetUserInfoByH5CodeRes, err error) {
	if _, err := oauthutil.VerifyApp(l.svcCtx.DB, req.AppID, req.AppSecret); err != nil {
		return nil, err
	}

	oauthCode, err := oauthutil.FindOAuthCode(l.svcCtx.DB, req.AppID, req.AuthCode)
	if err != nil {
		return nil, err
	}
	if err := oauthutil.ValidateOAuthCode(oauthCode); err != nil {
		return nil, err
	}
	if oauthCode.Scene != "h5_sso" {
		return nil, errors.New("????????")
	}

	if err := oauthutil.MarkOAuthCodeUsed(l.svcCtx.DB, oauthCode); err != nil {
		logx.Errorf("?? authCode ?????: %v", err)
		return nil, errors.New("???????")
	}

	userInfoRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: oauthCode.UserID,
	})
	if err != nil {
		logx.Errorf("????????: %v", err)
		return nil, errors.New("????????")
	}
	if userInfoRes.UserInfo == nil {
		return nil, errors.New("?????")
	}

	return &types.GetUserInfoByH5CodeRes{
		UserID:   userInfoRes.UserInfo.UserId,
		NickName: userInfoRes.UserInfo.NickName,
		Avatar:   userInfoRes.UserInfo.Avatar,
		Phone:    "",
		Email:    userInfoRes.UserInfo.Email,
	}, nil
}
