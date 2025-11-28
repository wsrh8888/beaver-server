package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackagesByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取表情包详情（用于数据同步）
func NewGetEmojiPackagesByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackagesByIdsLogic {
	return &GetEmojiPackagesByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiPackagesByIdsLogic) GetEmojiPackagesByIds(req *types.GetEmojiPackagesByIdsReq) (resp *types.GetEmojiPackagesByIdsRes, err error) {
	if len(req.Ids) == 0 {
		return &types.GetEmojiPackagesByIdsRes{
			Packages: make([]types.EmojiPackageDetailItem, 0),
		}, nil
	}

	// 根据UUID列表查询表情包详情
	var packages []emoji_models.EmojiPackage
	err = l.svcCtx.DB.Where("uuid IN ? AND status = ?", req.Ids, 1).Find(&packages).Error
	if err != nil {
		l.Errorf("查询表情包详情失败: uuids=%v, error=%v", req.Ids, err)
		return nil, err
	}

	l.Infof("批量查询表情包详情: 请求%d个, 返回%d个", len(req.Ids), len(packages))

	// 计算每个表情包的收藏数（可以考虑缓存这个统计数据）
	var packageItems []types.EmojiPackageDetailItem
	for _, pkg := range packages {
		// 查询收藏数
		var collectCount int64
		l.svcCtx.DB.Model(&emoji_models.EmojiPackageCollect{}).Where("package_id = ?", pkg.UUID).Count(&collectCount)

		packageItems = append(packageItems, types.EmojiPackageDetailItem{
			PackageID:    pkg.UUID,
			UUID:         pkg.UUID,
			Title:        pkg.Title,
			CoverFile:    pkg.CoverFile,
			UserID:       pkg.UserID,
			Description:  pkg.Description,
			Type:         pkg.Type,
			Status:       pkg.Status,
			CollectCount: int(collectCount),
			CreateAt:     time.Time(pkg.CreatedAt).UnixMilli(),
			UpdateAt:     time.Time(pkg.UpdatedAt).UnixMilli(),
			Version:      pkg.Version,
		})
	}

	return &types.GetEmojiPackagesByIdsRes{
		Packages: packageItems,
	}, nil
}
