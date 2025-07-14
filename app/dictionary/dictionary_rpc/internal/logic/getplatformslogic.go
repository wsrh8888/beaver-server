package logic

import (
	"context"

	"beaver/app/dictionary/dictionary_rpc/internal/svc"
	"beaver/app/dictionary/dictionary_rpc/types/dictionary_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPlatformsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPlatformsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPlatformsLogic {
	return &GetPlatformsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取平台列表
func (l *GetPlatformsLogic) GetPlatforms(in *dictionary_rpc.GetPlatformsReq) (*dictionary_rpc.GetPlatformsRes, error) {
	// 平台数据
	platforms := []*dictionary_rpc.PlatformInfo{
		{
			PlatformId:   1,
			PlatformName: "Windows",
			Description:  "Windows操作系统",
		},
		{
			PlatformId:   2,
			PlatformName: "MacOS",
			Description:  "苹果MacOS操作系统（电脑）",
		},
		{
			PlatformId:   3,
			PlatformName: "iOS",
			Description:  "苹果iOS操作系统（手机/平板）",
		},
		{
			PlatformId:   4,
			PlatformName: "Android",
			Description:  "安卓操作系统",
		},
		{
			PlatformId:   5,
			PlatformName: "HarmonyOS",
			Description:  "鸿蒙操作系统",
		},
	}

	return &dictionary_rpc.GetPlatformsRes{
		Platforms: platforms,
	}, nil
}
