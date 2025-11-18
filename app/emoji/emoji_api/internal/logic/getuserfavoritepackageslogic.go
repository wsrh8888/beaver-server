package logic

import (
	"context"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserFavoritePackagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserFavoritePackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFavoritePackagesLogic {
	return &GetUserFavoritePackagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserFavoritePackagesLogic) GetUserFavoritePackages(req *types.GetUserFavoritePackagesReq) (resp *types.GetUserFavoritePackagesRes, err error) {
	// 1. 查询用户收藏的表情包ID列表
	var packageCollects []emoji_models.EmojiPackageCollect
	err = l.svcCtx.DB.Where("user_id = ?", req.UserID).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Find(&packageCollects).Error

	if err != nil {
		logx.Errorf("查询用户收藏的表情包失败: %v", err)
		return nil, status.Error(codes.Internal, "查询收藏的表情包失败")
	}

	// 2. 获取收藏总数
	var total int64
	err = l.svcCtx.DB.Model(&emoji_models.EmojiPackageCollect{}).
		Where("user_id = ?", req.UserID).
		Count(&total).Error
	if err != nil {
		logx.Errorf("获取收藏总数失败: %v", err)
		return nil, status.Error(codes.Internal, "获取收藏总数失败")
	}

	// 如果没有收藏的表情包，直接返回空列表
	if len(packageCollects) == 0 {
		return &types.GetUserFavoritePackagesRes{
			Count: 0,
			List:  []types.EmojiPackageItem{},
		}, nil
	}

	// 3. 获取所有收藏的表情包ID
	packageIDs := make([]uint, len(packageCollects))
	for i, collect := range packageCollects {
		packageIDs[i] = collect.PackageID
	}

	// 4. 查询表情包详情
	var packages []emoji_models.EmojiPackage
	err = l.svcCtx.DB.Where("id IN ? AND status = ?", packageIDs, 1).Find(&packages).Error
	if err != nil {
		logx.Errorf("查询表情包详情失败: %v", err)
		return nil, status.Error(codes.Internal, "查询表情包详情失败")
	}

	// 创建表情包ID到对象的映射
	packageMap := make(map[uint]emoji_models.EmojiPackage)
	for _, p := range packages {
		packageMap[p.Id] = p
	}

	// 5. 获取每个表情包的表情数量
	emojiCounts := make(map[uint]int64)
	var emojiCountsData []struct {
		PackageID uint
		Count     int64
	}
	err = l.svcCtx.DB.Model(&emoji_models.EmojiPackageEmoji{}).
		Select("package_id, count(*) as count").
		Where("package_id IN ?", packageIDs).
		Group("package_id").
		Find(&emojiCountsData).Error
	if err != nil {
		logx.Errorf("获取表情数量失败: %v", err)
		return nil, status.Error(codes.Internal, "获取表情数量失败")
	}
	for _, c := range emojiCountsData {
		emojiCounts[c.PackageID] = c.Count
	}

	// 6. 获取每个表情包的收藏数
	collectCounts := make(map[uint]int64)
	var collectCountsData []struct {
		PackageID uint
		Count     int64
	}
	err = l.svcCtx.DB.Model(&emoji_models.EmojiPackageCollect{}).
		Select("package_id, count(*) as count").
		Where("package_id IN ?", packageIDs).
		Group("package_id").
		Find(&collectCountsData).Error
	if err != nil {
		logx.Errorf("获取收藏数失败: %v", err)
		return nil, status.Error(codes.Internal, "获取收藏数失败")
	}
	for _, c := range collectCountsData {
		collectCounts[c.PackageID] = c.Count
	}

	// 7. 构建返回数据
	packageItems := make([]types.EmojiPackageItem, 0, len(packageCollects))

	// 按照收藏的顺序构建响应
	for _, collect := range packageCollects {
		packageID := collect.PackageID
		package_, exists := packageMap[packageID]
		if !exists {
			continue // 跳过不存在或已禁用的表情包
		}

		packageItems = append(packageItems, types.EmojiPackageItem{
			PackageID:    package_.Id,
			Title:        package_.Title,
			CoverFile:    package_.CoverFile,
			Description:  package_.Description,
			Type:         package_.Type,
			CollectCount: int(collectCounts[packageID]),
			EmojiCount:   int(emojiCounts[packageID]),
			IsCollected:  true, // 这里一定是已收藏的
			IsAuthor:     package_.UserID == req.UserID,
		})
	}

	return &types.GetUserFavoritePackagesRes{
		Count: total,
		List:  packageItems,
	}, nil
}
