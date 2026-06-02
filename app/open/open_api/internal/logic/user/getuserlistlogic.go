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

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList(req *types.GetUserListReq) (resp *types.GetUserListRes, err error) {
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
	if len(req.UserIDs) == 0 {
		return nil, errors.New("userIds 不能为空")
	}
	for _, uid := range req.UserIDs {
		if uid != tokenRecord.UserID {
			return nil, errors.New("无权查询该用户信息")
		}
	}

	var users []user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id IN ?", req.UserIDs).Find(&users).Error; err != nil {
		return nil, err
	}

	scopes := oauthutil.ParseScopes(tokenRecord.Scope)
	userList := make([]types.GetUserListUserItem, 0, len(users))
	for _, user := range users {
		item := types.GetUserListUserItem{
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
		userList = append(userList, item)
	}

	return &types.GetUserListRes{Users: userList}, nil
}
