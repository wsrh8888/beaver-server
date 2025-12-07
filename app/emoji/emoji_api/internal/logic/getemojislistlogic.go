package logic

import (
	"context"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojisListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmojisListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojisListLogic {
	return &GetEmojisListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojisListLogic) GetEmojisList(req *types.GetEmojisListReq) (resp *types.GetEmojisListRes, err error) {
	// 获取用户收藏的表情ID列表
	var favoriteEmojis []emoji_models.EmojiCollectEmoji
	err = l.svcCtx.DB.Where("user_id = ?", req.UserID).Find(&favoriteEmojis).Error
	if err != nil {
		logx.Error("获取用户收藏的表情失败", err)
		return nil, err
	}

	// 提取表情UUID列表
	var emojiUUIDs []string
	for _, favorite := range favoriteEmojis {
		emojiUUIDs = append(emojiUUIDs, favorite.EmojiID)
	}

	if len(emojiUUIDs) == 0 {
		return &types.GetEmojisListRes{List: make([]types.EmojiItem, 0)}, nil
	}

	// 批量查询表情详情
	var emojis []emoji_models.Emoji
	err = l.svcCtx.DB.Where("uuid IN ? AND status = ?", emojiUUIDs, 1).Find(&emojis).Error
	if err != nil {
		logx.Error("获取表情详情失败", err)
		return nil, err
	}

	// 构建表情UUID到表情的映射
	emojiMap := make(map[string]emoji_models.Emoji)
	for _, emoji := range emojis {
		emojiMap[emoji.UUID] = emoji
	}

	// 构建响应数据
	var emojiItems []types.EmojiItem
	for _, favoriteEmoji := range favoriteEmojis {
		if emoji, exists := emojiMap[favoriteEmoji.EmojiID]; exists {
			emojiItems = append(emojiItems, types.EmojiItem{
				EmojiID:   emoji.UUID,
				FileKey:   emoji.FileKey, // 使用FileKey字段
				Title:     emoji.Title,
				PackageID: nil, // 在收藏表情列表中不显示包ID
			})
		}
	}

	resp = &types.GetEmojisListRes{
		List: emojiItems,
	}

	return resp, nil
}
