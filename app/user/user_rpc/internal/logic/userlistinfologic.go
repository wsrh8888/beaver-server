package logic

import (
	"context"

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
	// 对空数组进行处理
	if len(in.UserIdList) == 0 {
		return &user_rpc.UserListInfoRes{
			UserInfo: make(map[string]*user_rpc.UserInfo),
		}, nil
	}

	// 正确使用IN查询
	var userList []user_models.UserModel
	err := l.svcCtx.DB.Where("uuid IN ?", in.UserIdList).Find(&userList).Error

	if err != nil {
		l.Logger.Errorf("查询用户列表失败: %v", err)
		return nil, err
	}

	resp := &user_rpc.UserListInfoRes{
		UserInfo: make(map[string]*user_rpc.UserInfo, len(userList)),
	}

	for _, user := range userList {
		resp.UserInfo[user.UUID] = &user_rpc.UserInfo{
			NickName: user.NickName,
			FileName: user.FileName,
		}
	}

	l.Logger.Infof("查询到 %d 个用户信息，请求ID数量: %d", len(userList), len(in.UserIdList))
	return resp, nil
}
