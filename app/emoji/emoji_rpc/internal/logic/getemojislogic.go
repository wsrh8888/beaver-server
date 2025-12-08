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
	// 根据时间戳找到用户新增收藏的表情记录
	var collectRecords []struct {
		EmojiID   string    `gorm:"column:emoji_id"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}

	collectQuery := l.svcCtx.DB.Table("emoji_collect_emojis").Where("user_id = ?", in.UserId)

	// 时间戳过滤：只返回更新时间大于since的收藏记录
	if in.Since > 0 {
		sinceTime := time.UnixMilli(in.Since)
		collectQuery = collectQuery.Where("updated_at > ?", sinceTime)
	}

	err := collectQuery.Select("emoji_id, updated_at").Find(&collectRecords).Error
	if err != nil {
		l.Errorf("查询用户新增收藏表情记录失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	if len(collectRecords) == 0 {
		return &emoji_rpc.GetEmojisRes{
			EmojiVersions:   []*emoji_rpc.EmojiVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	// 提取表情ID列表
	emojiIDs := make([]string, 0, len(collectRecords))
	for _, record := range collectRecords {
		emojiIDs = append(emojiIDs, record.EmojiID)
	}

	// 获取这些表情的基础信息
	var emojis []struct {
		EmojiID string `gorm:"column:emoji_id"`
		Version int64  `gorm:"column:version"`
	}

	err = l.svcCtx.DB.Table("emojis").Where("emoji_id IN ? AND status = 1", emojiIDs).
		Select("emoji_id, version").Find(&emojis).Error
	if err != nil {
		l.Errorf("查询表情基础信息失败: emojiIds=%v, error=%v", emojiIDs, err)
		return nil, err
	}

	l.Infof("用户 %s 新增表情数据: 收藏记录=%d, 表情信息=%d", in.UserId, len(collectRecords), len(emojis))

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
