package logic

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RemoveEmojiFromPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 从表情包集合中移除表情图片
func NewRemoveEmojiFromPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveEmojiFromPackageLogic {
	return &RemoveEmojiFromPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveEmojiFromPackageLogic) RemoveEmojiFromPackage(req *types.RemoveEmojiFromPackageReq) (resp *types.RemoveEmojiFromPackageRes, err error) {
	// 转换PackageID为uint
	packageID, err := strconv.ParseUint(req.PackageID, 10, 32)
	if err != nil {
		logx.Errorf("表情包ID格式错误: %s", req.PackageID)
		return nil, errors.New("表情包ID格式错误")
	}

	// 转换EmojiID为uint
	emojiID, err := strconv.ParseUint(req.EmojiID, 10, 32)
	if err != nil {
		logx.Errorf("表情ID格式错误: %s", req.EmojiID)
		return nil, errors.New("表情ID格式错误")
	}

	// 检查表情包是否存在
	var pkg emoji_models.EmojiPackage
	err = l.svcCtx.DB.Where("id = ?", packageID).First(&pkg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("表情包不存在: %s", req.PackageID)
			return nil, errors.New("表情包不存在")
		}
		logx.Errorf("查询表情包失败: %v", err)
		return nil, errors.New("查询表情包失败")
	}

	// 检查表情是否存在
	var emoji emoji_models.Emoji
	err = l.svcCtx.DB.Where("id = ?", emojiID).First(&emoji).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("表情不存在: %s", req.EmojiID)
			return nil, errors.New("表情不存在")
		}
		logx.Errorf("查询表情失败: %v", err)
		return nil, errors.New("查询表情失败")
	}

	// 检查表情是否在表情包中
	var emojiPackageEmoji emoji_models.EmojiPackageEmoji
	err = l.svcCtx.DB.Where("package_id = ? AND emoji_id = ?", packageID, emojiID).First(&emojiPackageEmoji).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("表情不在该表情包中: %s", req.EmojiID)
			return nil, errors.New("表情不在该表情包中")
		}
		logx.Errorf("查询表情包关联失败: %v", err)
		return nil, errors.New("查询表情包关联失败")
	}

	// 删除表情包与表情的关联
	err = l.svcCtx.DB.Delete(&emojiPackageEmoji).Error
	if err != nil {
		logx.Errorf("从表情包移除表情失败: %v", err)
		return nil, errors.New("从表情包移除表情失败")
	}

	// 重新排序剩余的表情（可选，如果需要保持连续排序）
	// 这里可以选择重新排序或者保持现有排序

	return &types.RemoveEmojiFromPackageRes{}, nil
}
