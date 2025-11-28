package logic

import (
	"context"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GetEmojiPackageDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmojiPackageDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageDetailLogic {
	return &GetEmojiPackageDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiPackageDetailLogic) GetEmojiPackageDetail(req *types.GetEmojiPackageDetailReq) (*types.GetEmojiPackageDetailRes, error) {
	// 1. 获取表情包信息
	var emojiPackage emoji_models.EmojiPackage
	err := l.svcCtx.DB.Where("uuid = ?", req.PackageID).First(&emojiPackage).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "表情包不存在")
		}
		return nil, status.Error(codes.Internal, "获取表情包失败")
	}

	// 2. 检查表情包状态
	if emojiPackage.Status != 1 {
		return nil, status.Error(codes.PermissionDenied, "表情包已禁用")
	}

	// 3. 获取表情包中的表情列表
	// 首先获取关联关系
	var emojiPackageEmojis []emoji_models.EmojiPackageEmoji
	err = l.svcCtx.DB.Where("package_id = ?", req.PackageID).
		Order("sort_order").
		Find(&emojiPackageEmojis).Error

	if err != nil {
		logx.Errorf("查询表情包与表情关联关系失败: %v", err)
		return nil, status.Error(codes.Internal, "获取表情关联失败")
	}

	// 没有找到表情，返回空列表
	if len(emojiPackageEmojis) == 0 {
		return &types.GetEmojiPackageDetailRes{
			PackageID:    emojiPackage.UUID,
			Title:        emojiPackage.Title,
			CoverFile:    emojiPackage.CoverFile,
			Description:  emojiPackage.Description,
			Type:         emojiPackage.Type,
			CollectCount: 0,
			EmojiCount:   0,
			IsCollected:  false,
			IsAuthor:     emojiPackage.UserID == req.UserID,
			Emojis:       make([]types.EmojiItem, 0),
		}, nil
	}

	// 获取所有表情UUID
	emojiUUIDs := make([]string, len(emojiPackageEmojis))
	for i, emojiPackageEmoji := range emojiPackageEmojis {
		emojiUUIDs[i] = emojiPackageEmoji.EmojiID
	}

	// 查询表情详情
	var emojis []emoji_models.Emoji
	err = l.svcCtx.DB.Where("uuid IN ?", emojiUUIDs).Find(&emojis).Error
	if err != nil {
		logx.Errorf("查询表情详情失败: %v", err)
		return nil, status.Error(codes.Internal, "获取表情详情失败")
	}

	// 创建表情UUID到表情的映射，方便后续使用
	emojiMap := make(map[string]emoji_models.Emoji)
	for _, emoji := range emojis {
		emojiMap[emoji.UUID] = emoji
	}

	// 4. 获取收藏状态
	var collectCount int64
	err = l.svcCtx.DB.Model(&emoji_models.EmojiPackageCollect{}).Where("package_id = ?", req.PackageID).Count(&collectCount).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "获取收藏数失败")
	}

	// 5. 检查当前用户是否已收藏
	var isCollected bool
	err = l.svcCtx.DB.Model(&emoji_models.EmojiPackageCollect{}).
		Where("user_id = ? AND package_id = ? AND is_deleted = ?", req.UserID, req.PackageID, false).
		First(&emoji_models.EmojiPackageCollect{}).Error
	if err == nil {
		isCollected = true
	}

	// 6. 检查是否是作者
	isAuthor := emojiPackage.UserID == req.UserID

	// 7. 构建返回数据
	emojiItems := make([]types.EmojiItem, 0, len(emojiPackageEmojis))
	packageUUID := emojiPackage.UUID

	// 按照关联表中的顺序构建响应
	for _, emojiPackageEmoji := range emojiPackageEmojis {
		emoji, exists := emojiMap[emojiPackageEmoji.EmojiID]
		if !exists {
			continue // 跳过不存在的表情
		}

		emojiItems = append(emojiItems, types.EmojiItem{
			EmojiID:   emoji.UUID,
			FileName:  emoji.FileKey, // 使用FileKey字段
			Title:     emoji.Title,
			PackageID: &packageUUID,
		})
	}

	return &types.GetEmojiPackageDetailRes{
		PackageID:    emojiPackage.UUID,
		Title:        emojiPackage.Title,
		CoverFile:    emojiPackage.CoverFile,
		Description:  emojiPackage.Description,
		Type:         emojiPackage.Type,
		CollectCount: int(collectCount),
		EmojiCount:   len(emojiItems),
		IsCollected:  isCollected,
		IsAuthor:     isAuthor,
		Emojis:       emojiItems,
	}, nil
}
