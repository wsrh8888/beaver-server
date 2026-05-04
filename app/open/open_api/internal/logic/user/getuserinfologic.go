package user

import (
	"context"
	"errors"

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

// 获取用户信息
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (resp *types.GetUserInfoRes, err error) {
	// 1. 查询用户信息
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 返回用户信息
	return &types.GetUserInfoRes{
		User: types.UserInfo{
			UserID:   user.UserID,
			Nickname: user.NickName,
			Avatar:   user.Avatar,
			Phone:    user.Phone,
			Email:    user.Email,
		},
	}, nil
}
