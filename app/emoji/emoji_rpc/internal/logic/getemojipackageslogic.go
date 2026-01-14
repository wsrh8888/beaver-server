package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackagesLogic {
	return &GetEmojiPackagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiPackagesLogic) GetEmojiPackages(in *emoji_rpc.GetEmojiPackagesReq) (*emoji_rpc.GetEmojiPackagesRes, error) {
	// 1. 先查询用户收藏的表情包ID列表
	var collects []emoji_models.EmojiPackageCollect
	collectQuery := l.svcCtx.DB.Where("user_id = ? AND is_deleted = ?",
		in.UserId, false)

	// 增量同步：只返回更新时间大于since的收藏记录
	if in.Since > 0 {
		sinceTime := time.UnixMilli(in.Since)
		collectQuery = collectQuery.Where("updated_at > ?", sinceTime)
	}

	err := collectQuery.Find(&collects).Error
	if err != nil {
		l.Errorf("查询用户收藏表情包失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	// 如果没有收藏记录，直接返回空结果
	if len(collects) == 0 {
		l.Infof("用户 %s 没有收藏的表情包", in.UserId)
		return &emoji_rpc.GetEmojiPackagesRes{
			EmojiPackageVersions: []*emoji_rpc.EmojiPackageVersionItem{},
			ServerTimestamp:      time.Now().UnixMilli(),
		}, nil
	}

	// 2. 提取收藏的表情包ID列表
	var packageIDs []string
	for _, collect := range collects {
		packageIDs = append(packageIDs, collect.PackageID)
	}

	// 3. 查询这些表情包的详细信息
	var packages []emoji_models.EmojiPackage
	err = l.svcCtx.DB.Where("package_id IN (?) AND status = ?",
		packageIDs, 1).Find(&packages).Error // 1=正常状态
	if err != nil {
		l.Errorf("查询表情包详细信息失败: userId=%s, packageIds=%v, error=%v", in.UserId, packageIDs, err)
		return nil, err
	}

	l.Infof("查询到用户 %s 收藏的 %d 个表情包版本变更", in.UserId, len(packages))

	// 转换为版本摘要格式
	var packageVersions []*emoji_rpc.EmojiPackageVersionItem
	for _, pkg := range packages {
		packageVersions = append(packageVersions, &emoji_rpc.EmojiPackageVersionItem{
			PackageId: pkg.PackageID,
			Version:   pkg.Version,
		})
	}

	return &emoji_rpc.GetEmojiPackagesRes{
		EmojiPackageVersions: packageVersions,
		ServerTimestamp:      time.Now().UnixMilli(),
	}, nil
}
