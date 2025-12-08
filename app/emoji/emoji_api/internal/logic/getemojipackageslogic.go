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

type GetEmojiPackagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmojiPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackagesLogic {
	return &GetEmojiPackagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiPackagesLogic) GetEmojiPackages(req *types.GetEmojiPackagesReq) (*types.GetEmojiPackagesRes, error) {
	// 1. 构建查询条件
	query := l.svcCtx.DB.Model(&emoji_models.EmojiPackage{}).Where("status = ?", 1)

	// 2. 按分类筛选
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	// 3. 按类型筛选
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 4. 获取总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "获取总数失败")
	}

	// 5. 获取列表
	var packages []emoji_models.EmojiPackage
	err = query.Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&packages).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "获取列表失败")
	}

	// 6. 获取收藏状态和收藏数
	packageIDs := make([]string, len(packages))
	for i, p := range packages {
		packageIDs[i] = p.PackageID
	}

	// 获取收藏数
	collectCounts := make(map[string]int64)
	var collects []struct {
		PackageID string
		Count     int64
	}
	err = l.svcCtx.DB.Model(&emoji_models.EmojiPackageCollect{}).
		Select("package_id, count(*) as count").
		Where("package_id IN ?", packageIDs).
		Group("package_id").
		Find(&collects).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "获取收藏数失败")
	}
	for _, c := range collects {
		collectCounts[c.PackageID] = c.Count
	}

	// 获取表情数量
	emojiCounts := make(map[string]int64)
	var emojiCountsData []struct {
		PackageID string
		Count     int64
	}
	err = l.svcCtx.DB.Model(&emoji_models.EmojiPackageEmoji{}).
		Select("package_id, count(*) as count").
		Where("package_id IN ?", packageIDs).
		Group("package_id").
		Find(&emojiCountsData).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "获取表情数量失败")
	}
	for _, c := range emojiCountsData {
		emojiCounts[c.PackageID] = c.Count
	}

	// 获取当前用户的收藏状态
	userCollects := make(map[string]bool)
	if len(packageIDs) > 0 {
		var userCollectList []emoji_models.EmojiPackageCollect
		err = l.svcCtx.DB.Where("user_id = ? AND package_id IN ?", req.UserID, packageIDs).
			Find(&userCollectList).Error
		if err != nil {
			return nil, status.Error(codes.Internal, "获取收藏状态失败")
		}
		for _, c := range userCollectList {
			userCollects[c.PackageID] = true
		}
	}

	// 7. 构建返回数据
	list := make([]types.EmojiPackageItem, len(packages))
	for i, p := range packages {
		list[i] = types.EmojiPackageItem{
			PackageID:    p.PackageID,
			Title:        p.Title,
			CoverFile:    p.CoverFile,
			Description:  p.Description,
			Type:         p.Type,
			CollectCount: int(collectCounts[p.PackageID]),
			EmojiCount:   int(emojiCounts[p.PackageID]),
			IsCollected:  userCollects[p.PackageID],
			IsAuthor:     p.UserID == req.UserID,
		}
	}

	return &types.GetEmojiPackagesRes{
		Count: total,
		List:  list,
	}, nil
}
