package logic

import (
	"context"
	"errors"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
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
	// 生成表情版本号
	emojiVersion := l.svcCtx.VersionGen.GetNextVersion("emoji", "", "")
	if emojiVersion == -1 {
		logx.Error("生成表情版本号失败")
		return nil, errors.New("生成版本号失败")
	}

	// 创建表情
	emoji := emoji_models.Emoji{
		UUID:    uuid.New().String(),
		FileKey: req.FileKey, // 使用FileKey字段存储文件名
		Title:   req.Title,
		Version: emojiVersion,
	}

	// 保存到数据库
	err = l.svcCtx.DB.Create(&emoji).Error
	if err != nil {
		logx.Error("添加表情失败", err)
		return nil, err
	}

	// 如果指定了表情包ID，则添加到表情包
	if req.PackageID != "" {
		// 生成表情包表情关联的版本号（按表情包ID分区）
		packageEmojiVersion := l.svcCtx.VersionGen.GetNextVersion("emoji_package_emoji", "package_id", req.PackageID)
		if packageEmojiVersion == -1 {
			logx.Error("生成表情包表情关联版本号失败")
			return nil, errors.New("生成版本号失败")
		}

		// 创建表情包与表情的关联
		emojiPackageEmoji := emoji_models.EmojiPackageEmoji{
			UUID:      uuid.New().String(),
			PackageID: req.PackageID,
			EmojiID:   emoji.UUID,
			SortOrder: 0, // 默认排序
			Version:   packageEmojiVersion,
		}

		err = l.svcCtx.DB.Create(&emojiPackageEmoji).Error
		if err != nil {
			logx.Error("添加表情到表情包失败", err)
			return nil, err
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
		UUID:    uuid.New().String(),
		UserID:  req.UserID,
		EmojiID: emoji.UUID, // 使用新增表情的UUID
		Version: collectVersion,
	}

	err = l.svcCtx.DB.Create(&favoriteEmoji).Error
	if err != nil {
		logx.Error("收藏表情失败", err)
		return nil, err
	}

	return &types.AddEmojiRes{}, nil
}
