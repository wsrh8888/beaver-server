package logic

import (
	"context"

	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserVersionLogic {
	return &GetUserVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserVersionLogic) GetUserVersion(in *user_rpc.GetUserVersionReq) (*user_rpc.GetUserVersionRes, error) {
	// 获取用户模块的最新版本号
	version, err := l.svcCtx.VersionGen.GetCurrentVersion("users")
	if err != nil {
		l.Errorf("获取用户版本号失败: %v", err)
		return nil, err
	}

	l.Infof("获取用户版本号成功，用户ID: %s, 版本号: %d", in.UserId, version)

	return &user_rpc.GetUserVersionRes{
		LatestVersion: version,
	}, nil
}
