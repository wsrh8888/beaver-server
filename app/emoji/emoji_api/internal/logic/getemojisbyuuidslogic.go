package logic

import (
	"context"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojisByUuidsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmojisByUuidsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojisByUuidsLogic {
	return &GetEmojisByUuidsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojisByUuidsLogic) GetEmojisByUuids(req *types.GetEmojisByUuidsReq) (resp *types.GetEmojisByUuidsRes, err error) {
	if len(req.Uuids) == 0 {
		return &types.GetEmojisByUuidsRes{Emojis: []types.EmojiSimpleItem{}}, nil
	}

	var emojis []emoji_models.Emoji
	if err := l.svcCtx.DB.Where("uuid IN ? AND status = 1", req.Uuids).Find(&emojis).Error; err != nil {
		return nil, err
	}

	items := make([]types.EmojiSimpleItem, 0, len(emojis))
	for _, e := range emojis {
		items = append(items, types.EmojiSimpleItem{
			EmojiID: e.UUID,
			FileKey: e.FileKey,
			Title:   e.Title,
			Version: e.Version,
			Status:  e.Status,
		})
	}

	resp = &types.GetEmojisByUuidsRes{
		Emojis: items,
	}
	return resp, nil
}
