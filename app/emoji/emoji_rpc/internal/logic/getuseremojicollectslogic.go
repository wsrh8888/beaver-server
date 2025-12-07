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
	// 获取用户收藏的表情记录（返回收藏记录自身的版本，而不是表情版本）
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

	// 转换为收藏记录的版本摘要（使用收藏记录 UUID + 版本号）
	var emojiCollectVersions []*emoji_rpc.EmojiVersionItem
	for _, collect := range collectRecords {
		emojiCollectVersions = append(emojiCollectVersions, &emoji_rpc.EmojiVersionItem{
			Uuid:    collect.UUID,
			Version: collect.Version,
		})
	}

	return &emoji_rpc.GetUserEmojiCollectsRes{
		EmojiVersions:   emojiCollectVersions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
