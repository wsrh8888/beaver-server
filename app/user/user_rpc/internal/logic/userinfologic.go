package logic

import (
	"context"
	"encoding/json"
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

	err := l.svcCtx.DB.Take(&user, "user_id = ?", in.UserId).Error

	if err != nil {
		logx.Errorf("查询用户失败: %s", err.Error())
		return nil, errors.New("用户不存在")
	}

	byteData, _ := json.Marshal(user)

	return &user_rpc.UserInfoRes{
		Data: byteData,
	}, nil
}
