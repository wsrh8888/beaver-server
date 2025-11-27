package logic

import (
	"context"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojisListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojisListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojisListByIdsLogic {
	return &GetEmojisListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojisListByIdsLogic) GetEmojisListByIds(in *emoji_rpc.GetEmojisListByIdsReq) (*emoji_rpc.GetEmojisListByIdsRes, error) {
	if len(in.Ids) == 0 {
		return &emoji_rpc.GetEmojisListByIdsRes{Emojis: []*emoji_rpc.EmojiListById{}}, nil
	}

	var emojis []emoji_models.Emoji
	query := l.svcCtx.DB.Where("id IN (?)", in.Ids)

	// 时间戳过滤：只返回更新时间大于since的记录
	if in.Since > 0 {
		query = query.Where("updated_at > ?", in.Since)
	}

	err := query.Find(&emojis).Error
	if err != nil {
		l.Errorf("查询表情列表失败: ids=%v, since=%d, error=%v", in.Ids, in.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个表情详情", len(emojis))

	// 转换为响应格式
	var emojiList []*emoji_rpc.EmojiListById
	for _, emoji := range emojis {
		emojiList = append(emojiList, &emoji_rpc.EmojiListById{
			Id:       uint32(emoji.ID),
			FileKey:  emoji.FileKey,
			Title:    emoji.Title,
			AuthorId: emoji.AuthorID,
			Status:   int32(emoji.Status),
			Version:  emoji.Version,
			CreateAt: emoji.CreatedAt.UnixMilli(),
			UpdateAt: emoji.UpdatedAt.UnixMilli(),
		})
	}

	return &emoji_rpc.GetEmojisListByIdsRes{Emojis: emojiList}, nil
}
