package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojisByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取表情详情（用于数据同步）
func NewGetEmojisByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojisByIdsLogic {
	return &GetEmojisByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojisByIdsLogic) GetEmojisByIds(req *types.GetEmojisByIdsReq) (resp *types.GetEmojisByIdsRes, err error) {
	if len(req.Ids) == 0 {
		return &types.GetEmojisByIdsRes{
			Emojis: make([]types.EmojiDetailItem, 0),
		}, nil
	}

	// 根据UUID列表查询表情详情
	var emojis []emoji_models.Emoji
	err = l.svcCtx.DB.Where("uuid IN ? AND status = ?", req.Ids, 1).Find(&emojis).Error
	if err != nil {
		l.Errorf("查询表情详情失败: uuids=%v, error=%v", req.Ids, err)
		return nil, err
	}

	l.Infof("批量查询表情详情: 请求%d个, 返回%d个", len(req.Ids), len(emojis))

	// 获取表情UUID列表，用于查询关联的包信息
	emojiUUIDs := make([]string, len(emojis))
	for i, emoji := range emojis {
		emojiUUIDs[i] = emoji.UUID
	}

	// 查询表情包关联信息
	var packageEmojis []emoji_models.EmojiPackageEmoji
	if len(emojiUUIDs) > 0 {
		l.svcCtx.DB.Where("emoji_id IN ?", emojiUUIDs).Find(&packageEmojis)
	}

	// 建立表情UUID到包UUID的映射
	emojiToPackage := make(map[string]*string)
	for _, pe := range packageEmojis {
		if pe.PackageID != "" {
			emojiToPackage[pe.EmojiID] = &pe.PackageID
		}
	}

	// 转换为响应格式
	var emojiItems []types.EmojiDetailItem
	for _, emoji := range emojis {
		emojiItems = append(emojiItems, types.EmojiDetailItem{
			EmojiID:   emoji.UUID,
			UUID:      emoji.UUID,
			FileKey:   emoji.FileKey,
			Title:     emoji.Title,
			Status:    emoji.Status,
			Version:   emoji.Version,
			PackageID: emojiToPackage[emoji.UUID],
			CreateAt:  time.Time(emoji.CreatedAt).UnixMilli(),
			UpdateAt:  time.Time(emoji.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetEmojisByIdsRes{
		Emojis: emojiItems,
	}, nil
}
