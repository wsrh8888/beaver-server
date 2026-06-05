package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEmojisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEmojisLogic {
	return &ListEmojisLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListEmojisLogic) ListEmojis(in *emoji_rpc.ListEmojisReq) (*emoji_rpc.ListEmojisRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// package_id 非空：查包内表情，按 sort_order 排序
	if in.PackageId != "" {
		return l.listPackageEmojis(in.PackageId, page, pageSize)
	}

	db := l.svcCtx.DB.Model(&emoji_models.Emoji{})
	if in.EmojiId != "" {
		db = db.Where("emoji_id = ?", in.EmojiId)
	}
	if in.Title != "" {
		db = db.Where("title LIKE ?", "%"+in.Title+"%")
	}
	if in.StartTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", in.StartTime); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if in.EndTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", in.EndTime); err == nil {
			db = db.Where("created_at <= ?", t)
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []emoji_models.Emoji
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, err
	}

	items := make([]*emoji_rpc.EmojiItem, 0, len(list))
	for _, e := range list {
		items = append(items, toEmojiItem(e))
	}
	return &emoji_rpc.ListEmojisRes{Total: total, List: items}, nil
}

func (l *ListEmojisLogic) listPackageEmojis(packageID string, page, pageSize int) (*emoji_rpc.ListEmojisRes, error) {
	var relations []emoji_models.EmojiPackageEmoji
	if err := l.svcCtx.DB.Where("package_id = ?", packageID).
		Order("sort_order asc").Find(&relations).Error; err != nil {
		return nil, err
	}
	if len(relations) == 0 {
		return &emoji_rpc.ListEmojisRes{}, nil
	}

	emojiIDs := make([]string, len(relations))
	for i, r := range relations {
		emojiIDs[i] = r.EmojiID
	}

	var emojis []emoji_models.Emoji
	if err := l.svcCtx.DB.Where("emoji_id IN ?", emojiIDs).Find(&emojis).Error; err != nil {
		return nil, err
	}
	emojiMap := make(map[string]emoji_models.Emoji, len(emojis))
	for _, e := range emojis {
		emojiMap[e.EmojiID] = e
	}

	ordered := make([]emoji_models.Emoji, 0, len(relations))
	for _, r := range relations {
		if e, ok := emojiMap[r.EmojiID]; ok {
			ordered = append(ordered, e)
		}
	}

	total := int64(len(ordered))
	start := (page - 1) * pageSize
	if start >= len(ordered) {
		return &emoji_rpc.ListEmojisRes{Total: total}, nil
	}
	end := start + pageSize
	if end > len(ordered) {
		end = len(ordered)
	}

	items := make([]*emoji_rpc.EmojiItem, 0, end-start)
	for _, e := range ordered[start:end] {
		items = append(items, toEmojiItem(e))
	}
	return &emoji_rpc.ListEmojisRes{Total: total, List: items}, nil
}

func toEmojiItem(e emoji_models.Emoji) *emoji_rpc.EmojiItem {
	return &emoji_rpc.EmojiItem{
		EmojiId:   e.EmojiID,
		FileKey:   e.FileKey,
		Title:     e.Title,
		EmojiInfo: &emoji_rpc.EmojiInfoMsg{Width: int32(e.EmojiInfo.Width), Height: int32(e.EmojiInfo.Height)},
		Status:    int32(e.Status),
		CreatedAt: time.Time(e.CreatedAt).Format(time.RFC3339),
		UpdatedAt: time.Time(e.UpdatedAt).Format(time.RFC3339),
	}
}
