package user

import (
	"context"
	"errors"

	"beaver/app/open/constants"
	"beaver/app/open/open_api/internal/logic/oauthutil"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	user_models "beaver/app/user/user_models"

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

	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	scopes := oauthutil.ParseScopes(tokenRecord.Scope)
	item := types.GetUserInfoUserItem{
		UserID:   user.UserID,
		Nickname: user.NickName,
		Phone:    user.Phone,
	}
	if oauthutil.HasScope(scopes, constants.ScopeUserAvatarRead) {
		item.Avatar = user.Avatar
	}
	if oauthutil.HasScope(scopes, constants.ScopeUserEmailRead) {
		item.Email = user.Email
	}

	return &types.GetUserInfoRes{User: item}, nil
}
