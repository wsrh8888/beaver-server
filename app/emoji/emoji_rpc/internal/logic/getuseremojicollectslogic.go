package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserEmojiCollectsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserEmojiCollectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserEmojiCollectsLogic {
	return &GetUserEmojiCollectsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户收藏的表情版本摘要（用于datasync增量同步）
func (l *GetUserEmojiCollectsLogic) GetUserEmojiCollects(in *emoji_rpc.GetUserEmojiCollectsReq) (*emoji_rpc.GetUserEmojiCollectsRes, error) {
	// 首先获取用户收藏的表情ID列表
	var collectRecords []emoji_models.EmojiCollectEmoji
	collectQuery := l.svcCtx.DB.Where("user_id = ?", in.UserId)

	// 时间戳过滤：只返回更新时间大于since的收藏记录
	if in.Since > 0 {
		sinceTime := time.UnixMilli(in.Since)
		collectQuery = collectQuery.Where("updated_at > ?", sinceTime)
	}

	err := collectQuery.Find(&collectRecords).Error
	if err != nil {
		l.Errorf("查询用户收藏表情记录失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	if len(collectRecords) == 0 {
		return &emoji_rpc.GetUserEmojiCollectsRes{
			EmojiVersions:   []*emoji_rpc.EmojiVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	// 提取表情UUID列表
	emojiUUIDs := make([]string, 0, len(collectRecords))
	for _, record := range collectRecords {
		emojiUUIDs = append(emojiUUIDs, record.EmojiID)
	}

	// 根据表情UUID获取表情基础信息
	var emojis []emoji_models.Emoji
	err = l.svcCtx.DB.Where("uuid IN ?", emojiUUIDs).Find(&emojis).Error
	if err != nil {
		l.Errorf("查询表情基础信息失败: emojiUUIDs=%v, error=%v", emojiUUIDs, err)
		return nil, err
	}

	l.Infof("查询到用户 %s 的 %d 个收藏表情的基础信息", in.UserId, len(emojis))

	// 转换为版本摘要格式
	var emojiVersions []*emoji_rpc.EmojiVersionItem
	for _, emoji := range emojis {
		emojiVersions = append(emojiVersions, &emoji_rpc.EmojiVersionItem{
			Uuid:    emoji.UUID,
			Version: emoji.Version,
		})
	}

	return &emoji_rpc.GetUserEmojiCollectsRes{
		EmojiVersions:   emojiVersions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
