package logic

import (
	"context"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserVersionsLogic {
	return &UserVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserVersionsLogic) UserVersions(in *user_rpc.UserVersionsReq) (*user_rpc.UserVersionsRes, error) {
	// 对空数组进行处理
	if len(in.UserIds) == 0 {
		return &user_rpc.UserVersionsRes{
			UserVersions: make(map[string]int64),
		}, nil
	}

	// 查询指定用户ID列表的版本信息
	var userList []user_models.UserModel
	err := l.svcCtx.DB.Select("uuid, version").Where("uuid IN ?", in.UserIds).Find(&userList).Error

	if err != nil {
		l.Logger.Errorf("查询用户版本信息失败: %v", err)
		return nil, err
	}

	// 构造响应
	resp := &user_rpc.UserVersionsRes{
		UserVersions: make(map[string]int64, len(userList)),
	}

	for _, user := range userList {
		resp.UserVersions[user.UUID] = user.Version
	}

	// 为不存在的用户设置默认版本号0
	for _, userId := range in.UserIds {
		if _, exists := resp.UserVersions[userId]; !exists {
			resp.UserVersions[userId] = 0 // 用户不存在返回版本号0
		}
	}

	l.Logger.Infof("查询到 %d 个用户版本信息，请求ID数量: %d", len(userList), len(in.UserIds))
	return resp, nil
}
