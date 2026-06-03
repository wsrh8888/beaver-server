package user

import (
	"context"
	"errors"

	"beaver/app/open/constants"
	"beaver/app/open/open_api/internal/logic/oauthutil"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (resp *types.GetUserInfoRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("userId 不能为空")
	}

	tokenRecord, err := oauthutil.ValidateAccessTokenWithScopes(
		l.svcCtx.DB,
		req.Authorization,
		constants.ScopeUserProfileRead,
	)
	if err != nil {
		return nil, err
	}
	if err := oauthutil.RequireUserToken(tokenRecord); err != nil {
		return nil, err
	}
	if req.UserID != tokenRecord.UserID {
		return nil, errors.New("无权查询该用户信息")
	}

	userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: req.UserID})
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	scopes := oauthutil.ParseScopes(tokenRecord.Scope)
	item := types.GetUserInfoUserItem{
		UserID:   userRes.UserInfo.UserId,
		Nickname: userRes.UserInfo.NickName,
	}
	if oauthutil.HasScope(scopes, constants.ScopeUserAvatarRead) {
		item.Avatar = userRes.UserInfo.Avatar
	}
	if oauthutil.HasScope(scopes, constants.ScopeUserEmailRead) {
		item.Email = userRes.UserInfo.Email
	}

	return &types.GetUserInfoRes{User: item}, nil
}
