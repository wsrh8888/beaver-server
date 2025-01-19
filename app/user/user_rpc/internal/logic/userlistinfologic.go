package logic

import (
	"context"
	"fmt"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserListInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListInfoLogic {
	return &UserListInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserListInfoLogic) UserListInfo(in *user_rpc.UserListInfoReq) (*user_rpc.UserListInfoRes, error) {
	var userList []user_models.UserModel

	l.svcCtx.DB.Find(&userList, "uuid = ?", in.UserIdList)

	resp := new(user_rpc.UserListInfoRes)
	resp.UserInfo = make(map[string]*user_rpc.UserInfo, 0)

	for _, i2 := range userList {
		resp.UserInfo[i2.UUID] = &user_rpc.UserInfo{
			NickName: i2.NickName,
			Avatar:   i2.Avatar,
		}
	}
	fmt.Println(userList, in.UserIdList)

	return resp, nil
}
