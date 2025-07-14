package logic

import (
	"context"

	"beaver/app/dictionary/dictionary_rpc/internal/svc"
	"beaver/app/dictionary/dictionary_rpc/types/dictionary_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArchitecturesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArchitecturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArchitecturesLogic {
	return &GetArchitecturesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取架构列表
func (l *GetArchitecturesLogic) GetArchitectures(in *dictionary_rpc.GetArchitecturesReq) (*dictionary_rpc.GetArchitecturesRes, error) {
	// 架构数据
	architectures := []*dictionary_rpc.ArchitectureInfo{
		// H5架构
		{
			ArchId:      0,
			ArchName:    "H5",
			Description: "H5网页版本",
			PlatformId:  0, // H5不属于特定平台
		},
		// Windows架构
		{
			ArchId:      1,
			ArchName:    "WinX64",
			Description: "Windows Intel x64 (64位)",
			PlatformId:  1,
		},
		{
			ArchId:      2,
			ArchName:    "WinArm64",
			Description: "Windows ARM64 (Surface等设备)",
			PlatformId:  1,
		},
		// MacOS架构
		{
			ArchId:      3,
			ArchName:    "MacIntel",
			Description: "MacOS Intel版本",
			PlatformId:  2,
		},
		{
			ArchId:      4,
			ArchName:    "MacApple",
			Description: "MacOS Apple Silicon (M1/M2/M3系列)",
			PlatformId:  2,
		},
		// iOS架构
		{
			ArchId:      5,
			ArchName:    "iOS",
			Description: "iOS通用版本",
			PlatformId:  3,
		},
		// Android架构
		{
			ArchId:      6,
			ArchName:    "Android",
			Description: "Android通用版本",
			PlatformId:  4,
		},
		// 鸿蒙架构
		{
			ArchId:      7,
			ArchName:    "HarmonyOS",
			Description: "鸿蒙系统通用版本",
			PlatformId:  5,
		},
	}

	// 如果指定了平台ID，则过滤
	if in.PlatformId > 0 {
		filtered := make([]*dictionary_rpc.ArchitectureInfo, 0)
		for _, arch := range architectures {
			if arch.PlatformId == in.PlatformId {
				filtered = append(filtered, arch)
			}
		}
		architectures = filtered
	}

	return &dictionary_rpc.GetArchitecturesRes{
		Architectures: architectures,
	}, nil
}
