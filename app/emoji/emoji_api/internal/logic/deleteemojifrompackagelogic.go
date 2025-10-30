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

type DeleteEmojiFromPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteEmojiFromPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEmojiFromPackageLogic {
	return &DeleteEmojiFromPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteEmojiFromPackageLogic) DeleteEmojiFromPackage(req *types.DeleteEmojiFromPackageReq) (*types.DeleteEmojiFromPackageRes, error) {
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
		return nil, status.Error(codes.PermissionDenied, "只有表情包作者可以删除表情")
	}

	// 4. 检查表情是否存在
	var emoji emoji_models.Emoji
	err = l.svcCtx.DB.First(&emoji, req.EmojiID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "表情不存在")
		}
		return nil, status.Error(codes.Internal, "获取表情失败")
	}

	// 5. 检查表情和表情包的关联关系是否存在
	var emojiPackageEmoji emoji_models.EmojiPackageEmoji
	err = l.svcCtx.DB.Where("package_id = ? AND emoji_id = ?", req.PackageID, req.EmojiID).First(&emojiPackageEmoji).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "表情不在该表情包中")
		}
		return nil, status.Error(codes.Internal, "查询关联关系失败")
	}

	// 6. 开启事务
	tx := l.svcCtx.DB.Begin()

	// 7. 删除表情和表情包的关联关系
	err = tx.Where("package_id = ? AND emoji_id = ?", req.PackageID, req.EmojiID).Delete(&emoji_models.EmojiPackageEmoji{}).Error
	if err != nil {
		tx.Rollback()
		logx.Error("删除表情和表情包关联关系失败", err)
		return nil, status.Error(codes.Internal, "删除表情失败")
	}

	// 8. 检查该表情是否还被其他表情包使用
	var count int64
	err = tx.Model(&emoji_models.EmojiPackageEmoji{}).Where("emoji_id = ?", req.EmojiID).Count(&count).Error
	if err != nil {
		tx.Rollback()
		logx.Error("查询表情使用情况失败", err)
		return nil, status.Error(codes.Internal, "删除表情失败")
	}

	// 9. 如果表情不再被任何表情包使用，考虑是否删除表情本身（可选）
	// 这里选择保留表情本身，因为它可能被用户收藏

	// 10. 提交事务
	tx.Commit()

	return &types.DeleteEmojiFromPackageRes{}, nil
}
