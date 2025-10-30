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
	var favoriteEmojis []emoji_models.EmojiCollectEmoji
	err = l.svcCtx.DB.Where("user_id = ?", req.UserID).Preload("EmojiModel").Find(&favoriteEmojis).Error
	if err != nil {
		logx.Error("获取用户收藏的表情失败", err)
		return nil, err
	}

	var emojiItems []types.EmojiItem
	for _, favoriteEmoji := range favoriteEmojis {
		emojiItems = append(emojiItems, types.EmojiItem{
			EmojiID:   favoriteEmoji.EmojiModel.Id,
			FileName:  favoriteEmoji.EmojiModel.FileName,
			Title:     favoriteEmoji.EmojiModel.Title,
			PackageID: nil, // 在收藏表情列表中不显示包ID
		})
	}

	resp = &types.GetEmojisListRes{
		List: emojiItems,
	}

	return resp, nil
}
