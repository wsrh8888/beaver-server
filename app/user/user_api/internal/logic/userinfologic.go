package logic

import (
	"context"
	"encoding/json"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoRes, err error) {
	res, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: req.UserID})
	if err != nil {
		return nil, err
	}
	var user user_models.UserModel

	err = json.Unmarshal(res.Data, &user)
	if err != nil {
		return nil, err
	}
	resp = &types.UserInfoRes{
		UserID:   user.UUID,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Abstract: user.Abstract,
		Phone:    user.Phone,
	}

	return
}
