package logic

import (
	"context"
	"errors"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AddEmojiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddEmojiLogic {
	return &AddEmojiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddEmojiLogic) AddEmoji(req *types.AddEmojiReq) (resp *types.AddEmojiRes, err error) {
	// 先按 FileKey 查重，已有则复用，不重复落库
	var emoji emoji_models.Emoji
	err = l.svcCtx.DB.Where("file_key = ?", req.FileKey).First(&emoji).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 仅当不存在时创建新的 emoji
	if errors.Is(err, gorm.ErrRecordNotFound) {
		emojiVersion := l.svcCtx.VersionGen.GetNextVersion("emoji", "", "")
		if emojiVersion == -1 {
			logx.Error("生成表情版本号失败")
			return nil, errors.New("生成版本号失败")
		}

		emoji = emoji_models.Emoji{
			EmojiID: uuid.New().String(),
			FileKey: req.FileKey,
			Title:   req.Title,
			Version: emojiVersion,
		}

		if err := l.svcCtx.DB.Create(&emoji).Error; err != nil {
			logx.Error("添加表情失败", err)
			return nil, err
		}
	}

	// 如果指定了表情包ID，则添加到表情包
	if req.PackageID != "" {
		var existPkgEmoji emoji_models.EmojiPackageEmoji
		err = l.svcCtx.DB.Where("package_id = ? AND emoji_id = ?", req.PackageID, emoji.EmojiID).First(&existPkgEmoji).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			packageEmojiVersion := l.svcCtx.VersionGen.GetNextVersion("emoji_package_emoji", "package_id", req.PackageID)
			if packageEmojiVersion == -1 {
				logx.Error("生成表情包表情关联版本号失败")
				return nil, errors.New("生成版本号失败")
			}

			emojiPackageEmoji := emoji_models.EmojiPackageEmoji{
				RelationID: uuid.New().String(),
				PackageID:  req.PackageID,
				EmojiID:    emoji.EmojiID,
				SortOrder:  0,
				Version:    packageEmojiVersion,
			}

			if err := l.svcCtx.DB.Create(&emojiPackageEmoji).Error; err != nil {
				logx.Error("添加表情到表情包失败", err)
				return nil, err
			}
		}
	}

	// 生成收藏版本号（按用户ID分区）
	collectVersion := l.svcCtx.VersionGen.GetNextVersion("emoji_collect", "user_id", req.UserID)
	if collectVersion == -1 {
		logx.Error("生成收藏版本号失败")
		return nil, errors.New("生成版本号失败")
	}

	// 添加表情并收藏
	favoriteEmoji := emoji_models.EmojiCollectEmoji{
		EmojiCollectID: uuid.New().String(),
		UserID:         req.UserID,
		EmojiID:        emoji.EmojiID,
		Version:        collectVersion,
	}

	// 去重：同一用户对同一 emoji 已收藏则跳过创建
	var existFavorite emoji_models.EmojiCollectEmoji
	err = l.svcCtx.DB.Where("user_id = ? AND emoji_id = ?", req.UserID, emoji.EmojiID).First(&existFavorite).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := l.svcCtx.DB.Create(&favoriteEmoji).Error; err != nil {
			logx.Error("收藏表情失败", err)
			return nil, err
		}
	}

	return &types.AddEmojiRes{}, nil
}
