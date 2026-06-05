package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEmojiCollectsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListEmojiCollectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEmojiCollectsLogic {
	return &ListEmojiCollectsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListEmojiCollectsLogic) ListEmojiCollects(in *emoji_rpc.ListEmojiCollectsReq) (*emoji_rpc.ListEmojiCollectsRes, error) {
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

	db := l.svcCtx.DB.Model(&emoji_models.EmojiCollectEmoji{})
	if in.UserId != "" {
		db = db.Where("user_id = ?", in.UserId)
	}
	if in.EmojiId != "" {
		db = db.Where("emoji_id = ?", in.EmojiId)
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

	var collects []emoji_models.EmojiCollectEmoji
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&collects).Error; err != nil {
		return nil, err
	}

	emojiIDs := make([]string, 0, len(collects))
	for _, c := range collects {
		emojiIDs = append(emojiIDs, c.EmojiID)
	}
	emojiMap := map[string]emoji_models.Emoji{}
	if len(emojiIDs) > 0 {
		var emojis []emoji_models.Emoji
		if err := l.svcCtx.DB.Where("emoji_id IN ?", emojiIDs).Find(&emojis).Error; err == nil {
			for _, e := range emojis {
				emojiMap[e.EmojiID] = e
			}
		}
	}

	items := make([]*emoji_rpc.EmojiCollectItem, 0, len(collects))
	for _, c := range collects {
		title, fileKey := "", ""
		if e, ok := emojiMap[c.EmojiID]; ok {
			title, fileKey = e.Title, e.FileKey
		}
		items = append(items, &emoji_rpc.EmojiCollectItem{
			CollectId:    c.EmojiCollectID,
			UserId:       c.UserID,
			EmojiId:      c.EmojiID,
			EmojiTitle:   title,
			EmojiFileKey: fileKey,
			CreatedAt:    time.Time(c.CreatedAt).Format(time.RFC3339),
			UpdatedAt:    time.Time(c.UpdatedAt).Format(time.RFC3339),
		})
	}
	return &emoji_rpc.ListEmojiCollectsRes{Total: total, List: items}, nil
}
