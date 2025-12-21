package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojisLogic {
	return &GetEmojisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ========== 用户相关表情基础数据同步 ==========
func (l *GetEmojisLogic) GetEmojis(in *emoji_rpc.GetEmojisReq) (*emoji_rpc.GetEmojisRes, error) {
	// 获取用户相关的表情ID列表：
	// 1. 用户直接收藏的表情
	// 2. 用户收藏的表情包中包含的表情

	var userEmojiIDs []string
	var packageEmojiIDs []string

	// 1. 查询用户收藏的表情ID
	var userCollectRecords []struct {
		EmojiID string `gorm:"column:emoji_id"`
	}
	err := l.svcCtx.DB.Table("emoji_collect_emojis").Where("user_id = ?", in.UserId).
		Select("emoji_id").Find(&userCollectRecords).Error
	if err != nil {
		l.Errorf("查询用户收藏表情失败: userId=%s, error=%v", in.UserId, err)
		return nil, err
	}

	for _, record := range userCollectRecords {
		userEmojiIDs = append(userEmojiIDs, record.EmojiID)
	}

	// 2. 查询用户收藏的表情包中包含的表情ID
	var packageCollectRecords []struct {
		PackageID string `gorm:"column:package_id"`
	}
	err = l.svcCtx.DB.Table("emoji_package_collects").Where("user_id = ?", in.UserId).
		Select("package_id").Find(&packageCollectRecords).Error
	if err != nil {
		l.Errorf("查询用户收藏表情包失败: userId=%s, error=%v", in.UserId, err)
		return nil, err
	}

	if len(packageCollectRecords) > 0 {
		packageIDs := make([]string, 0, len(packageCollectRecords))
		for _, record := range packageCollectRecords {
			packageIDs = append(packageIDs, record.PackageID)
		}

		// 查询这些表情包包含的所有表情ID
		var packageContentRecords []struct {
			EmojiID string `gorm:"column:emoji_id"`
		}
		err = l.svcCtx.DB.Table("emoji_package_emojis").Where("package_id IN ?", packageIDs).
			Select("emoji_id").Find(&packageContentRecords).Error
		if err != nil {
			l.Errorf("查询表情包内容失败: userId=%s, packageIds=%v, error=%v", in.UserId, packageIDs, err)
			return nil, err
		}

		for _, record := range packageContentRecords {
			packageEmojiIDs = append(packageEmojiIDs, record.EmojiID)
		}
	}

	// 合并并去重表情ID
	emojiIDMap := make(map[string]bool)
	for _, id := range userEmojiIDs {
		emojiIDMap[id] = true
	}
	for _, id := range packageEmojiIDs {
		emojiIDMap[id] = true
	}

	if len(emojiIDMap) == 0 {
		return &emoji_rpc.GetEmojisRes{
			EmojiVersions:   []*emoji_rpc.EmojiVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	// 转换为ID列表
	allEmojiIDs := make([]string, 0, len(emojiIDMap))
	for id := range emojiIDMap {
		allEmojiIDs = append(allEmojiIDs, id)
	}

	// 获取这些表情的基础信息
	var emojis []struct {
		EmojiID string `gorm:"column:emoji_id"`
		Version int64  `gorm:"column:version"`
	}

	err = l.svcCtx.DB.Table("emojis").Where("emoji_id IN ? AND status = 1", allEmojiIDs).
		Select("emoji_id, version").Find(&emojis).Error
	if err != nil {
		l.Errorf("查询表情基础信息失败: userId=%s, emojiIds=%v, error=%v", in.UserId, allEmojiIDs, err)
		return nil, err
	}

	l.Infof("用户 %s 表情基础数据同步: 直接收藏=%d, 表情包包含=%d, 总计=%d, 有效表情=%d",
		in.UserId, len(userEmojiIDs), len(packageEmojiIDs), len(emojiIDMap), len(emojis))

	// 转换为版本摘要格式
	var emojiVersions []*emoji_rpc.EmojiVersionItem
	for _, emoji := range emojis {
		emojiVersions = append(emojiVersions, &emoji_rpc.EmojiVersionItem{
			EmojiId: emoji.EmojiID,
			Version: emoji.Version,
		})
	}

	return &emoji_rpc.GetEmojisRes{
		EmojiVersions:   emojiVersions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
