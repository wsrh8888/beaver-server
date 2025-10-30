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

type BatchAddEmojiToPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchAddEmojiToPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchAddEmojiToPackageLogic {
	return &BatchAddEmojiToPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchAddEmojiToPackageLogic) BatchAddEmojiToPackage(req *types.BatchAddEmojiToPackageReq) (*types.BatchAddEmojiToPackageRes, error) {
	// 1. 检查表情包是否存在
	var emojiPackage emoji_models.EmojiPackage
	err := l.svcCtx.DB.First(&emojiPackage, req.PackageID).Error
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

	// 3. 检查是否是表情包作者
	if emojiPackage.UserID != req.UserID {
		return nil, status.Error(codes.PermissionDenied, "只有表情包作者可以添加表情")
	}

	// 4. 开启事务
	tx := l.svcCtx.DB.Begin()

	// 5. 批量创建表情
	emojis := make([]emoji_models.Emoji, len(req.Emojis))
	for i, emojiReq := range req.Emojis {
		emojis[i] = emoji_models.Emoji{
			FileName: emojiReq.FileName,
			Title:    emojiReq.Title,
			AuthorID: req.UserID,
		}
	}

	// 6. 保存表情到数据库
	err = tx.Create(&emojis).Error
	if err != nil {
		tx.Rollback()
		logx.Error("批量添加表情失败", err)
		return nil, status.Error(codes.Internal, "批量添加表情失败")
	}

	// 7. 创建表情和表情包的关联关系
	emojiPackageEmojis := make([]emoji_models.EmojiPackageEmoji, len(emojis))
	for i, emoji := range emojis {
		emojiPackageEmojis[i] = emoji_models.EmojiPackageEmoji{
			PackageID: req.PackageID,
			EmojiID:   emoji.Id,
			SortOrder: i, // 使用索引作为排序顺序
		}
	}

	// 8. 保存关联关系到数据库
	err = tx.Create(&emojiPackageEmojis).Error
	if err != nil {
		tx.Rollback()
		logx.Error("批量创建表情关联关系失败", err)
		return nil, status.Error(codes.Internal, "批量创建表情关联关系失败")
	}

	// 9. 提交事务
	tx.Commit()

	// 10. 收集新创建的表情ID
	emojiIDs := make([]uint, len(emojis))
	for i, emoji := range emojis {
		emojiIDs[i] = emoji.Id
	}

	return &types.BatchAddEmojiToPackageRes{
		EmojiIDs: emojiIDs,
	}, nil
}
