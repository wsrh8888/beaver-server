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

	// 构建查询条件
	query := l.svcCtx.DB.Model(&user_models.UserModel{}).Where("uuid IN ?", in.UserIdList)

	// 如果提供了时间戳，则只返回该时间之后更新的用户
	if in.SinceTimestamp > 0 {
		query = query.Where("updated_at > ?", in.SinceTimestamp)
	}

	var userList []user_models.UserModel
	err := query.Find(&userList).Error

	if err != nil {
		l.Logger.Errorf("查询用户列表失败: %v", err)
		return nil, err
	}

	resp := &user_rpc.UserListInfoRes{
		UserInfo: make(map[string]*user_rpc.UserInfo, len(userList)),
	}

	for _, user := range userList {
		resp.UserInfo[user.UUID] = &user_rpc.UserInfo{
			UserId:   user.UUID, // 保持向后兼容
			NickName: user.NickName,
			Avatar:   user.Avatar,
			Version:  user.Version,
			Email:    user.Email,
			Abstract: user.Abstract,
		}
	}

	if in.SinceTimestamp > 0 {
		l.Logger.Infof("增量查询用户，时间戳: %d，查询到 %d 个用户信息，请求ID数量: %d",
			in.SinceTimestamp, len(userList), len(in.UserIdList))
	} else {
		l.Logger.Infof("全量查询用户，查询到 %d 个用户信息，请求ID数量: %d",
			len(userList), len(in.UserIdList))
	}

	return resp, nil
}
