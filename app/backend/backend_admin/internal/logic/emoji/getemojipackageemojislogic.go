package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetEmojiPackageEmojisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取表情包内的表情图片列表
func NewGetEmojiPackageEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageEmojisLogic {
	return &GetEmojiPackageEmojisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiPackageEmojisLogic) GetEmojiPackageEmojis(req *types.GetEmojiPackageEmojisReq) (resp *types.GetEmojiPackageEmojisRes, err error) {
	// 检查表情包是否存在
	var pkg emoji_models.EmojiPackage
	err = l.svcCtx.DB.Where("uuid = ?", req.PackageUUID).First(&pkg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("表情包不存在: %s", req.PackageUUID)
			return nil, errors.New("表情包不存在")
		}
		logx.Errorf("查询表情包失败: %v", err)
		return nil, errors.New("查询表情包失败")
	}

	// 先查询关联关系
	var emojiPackageEmojis []emoji_models.EmojiPackageEmoji
	err = l.svcCtx.DB.Where("package_id = ?", pkg.Id).
		Order("sort_order asc").
		Find(&emojiPackageEmojis).Error
	if err != nil {
		logx.Errorf("查询表情包关联失败: %v", err)
		return nil, err
	}

	// 如果没有表情，返回空列表
	if len(emojiPackageEmojis) == 0 {
		return &types.GetEmojiPackageEmojisRes{
			List:  []types.GetEmojiPackageEmojisItem{},
			Total: 0,
		}, nil
	}

	// 获取所有表情ID
	emojiIDs := make([]uint, len(emojiPackageEmojis))
	for i, emojiPackageEmoji := range emojiPackageEmojis {
		emojiIDs[i] = emojiPackageEmoji.EmojiID
	}

	// 查询表情详情
	var emojis []emoji_models.Emoji
	err = l.svcCtx.DB.Where("id IN ?", emojiIDs).Find(&emojis).Error
	if err != nil {
		logx.Errorf("查询表情详情失败: %v", err)
		return nil, err
	}

	// 创建表情UUID到表情的映射
	emojiMap := make(map[string]emoji_models.Emoji)
	for _, emoji := range emojis {
		emojiMap[emoji.UUID] = emoji
	}

	// 按照关联表中的顺序构建结果
	var orderedEmojis []emoji_models.Emoji
	for _, emojiPackageEmoji := range emojiPackageEmojis {
		if emoji, exists := emojiMap[emojiPackageEmoji.EmojiID]; exists {
			orderedEmojis = append(orderedEmojis, emoji)
		}
	}

	// 分页参数校验
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 手动分页
	total := int64(len(orderedEmojis))
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(orderedEmojis) {
		start = len(orderedEmojis)
	}
	if end > len(orderedEmojis) {
		end = len(orderedEmojis)
	}

	pagedEmojis := orderedEmojis[start:end]

	// 转换为响应格式
	var list []types.GetEmojiPackageEmojisItem
	for _, emoji := range pagedEmojis {
		list = append(list, types.GetEmojiPackageEmojisItem{
			UUID:       emoji.UUID,
			FileKey:    emoji.FileKey,
			Title:      emoji.Title,
			AuthorID:   "", // 暂时为空，后续可从其他途径获取
			CreateTime: emoji.CreatedAt.String(),
			UpdateTime: emoji.UpdatedAt.String(),
		})
	}

	return &types.GetEmojiPackageEmojisRes{
		List:  list,
		Total: total,
	}, nil
}
