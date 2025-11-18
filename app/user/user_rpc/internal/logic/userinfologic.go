package logic

import (
	"context"
	"errors"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user_rpc.UserInfoReq) (*user_rpc.UserInfoRes, error) {
	var user user_models.UserModel

	err := l.svcCtx.DB.Take(&user, "uuid = ?", in.UserID).Error

	if err != nil {
		logx.Errorf("查询用户失败: %s", err.Error())
		return nil, errors.New("用户不存在")
	}

	return &user_rpc.UserInfoRes{
		UserInfo: &user_rpc.UserInfo{
			UserId:   user.UUID,
			NickName: user.NickName,
			Avatar:   user.Avatar,
			Version:  user.Version,
			Email:    user.Email,
			Abstract: user.Abstract,
		},
	}, nil
}
