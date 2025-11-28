package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageContentsByPackageIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取表情包内容详情（同步用）
func NewGetEmojiPackageContentsByPackageIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageContentsByPackageIdsLogic {
	return &GetEmojiPackageContentsByPackageIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiPackageContentsByPackageIdsLogic) GetEmojiPackageContentsByPackageIds(req *types.GetEmojiPackageContentsByPackageIdsReq) (resp *types.GetEmojiPackageContentsByPackageIdsRes, err error) {
	if len(req.PackageIds) == 0 {
		return &types.GetEmojiPackageContentsByPackageIdsRes{
			Contents: make([]types.EmojiPackageContentDetailItem, 0),
		}, nil
	}

	// 根据表情包ID列表查询内容详情
	var contents []emoji_models.EmojiPackageEmoji
	err = l.svcCtx.DB.Where("package_id IN ?", req.PackageIds).Find(&contents).Error
	if err != nil {
		l.Errorf("查询表情包内容详情失败: packageIds=%v, error=%v", req.PackageIds, err)
		return nil, err
	}

	l.Infof("批量查询表情包内容详情: 请求%d个包, 返回%d条内容", len(req.PackageIds), len(contents))

	// 转换为响应格式
	var contentItems []types.EmojiPackageContentDetailItem
	for _, content := range contents {
		contentItems = append(contentItems, types.EmojiPackageContentDetailItem{
			PackageID: content.PackageID,
			EmojiID:   content.EmojiID,
			SortOrder: content.SortOrder,
			Version:   content.Version,
			CreateAt:  time.Time(content.CreatedAt).UnixMilli(),
			UpdateAt:  time.Time(content.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetEmojiPackageContentsByPackageIdsRes{
		Contents: contentItems,
	}, nil
}
