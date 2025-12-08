package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiCollectsByUuidsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按收藏ID批量获取表情收藏记录（同步补齐）
func NewGetEmojiCollectsByUuidsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiCollectsByUuidsLogic {
	return &GetEmojiCollectsByUuidsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiCollectsByUuidsLogic) GetEmojiCollectsByUuids(req *types.GetEmojiCollectsByUuidsReq) (resp *types.GetEmojiCollectsByUuidsRes, err error) {
	if len(req.Ids) == 0 {
		return &types.GetEmojiCollectsByUuidsRes{
			Collects: make([]types.EmojiCollectDetailItem, 0),
		}, nil
	}

	var collects []emoji_models.EmojiCollectEmoji
	err = l.svcCtx.DB.Where("emoji_collect_id IN ?", req.Ids).Find(&collects).Error
	if err != nil {
		l.Errorf("按收藏ID批量查询表情收藏记录失败: ids=%v, error=%v", req.Ids, err)
		return nil, err
	}

	collectItems := make([]types.EmojiCollectDetailItem, 0, len(collects))
	for _, collect := range collects {
		collectItems = append(collectItems, types.EmojiCollectDetailItem{
			EmojiCollectID: collect.EmojiCollectID,
			UserID:         collect.UserID,
			EmojiID:        collect.EmojiID,
			IsDeleted:      collect.IsDeleted,
			Version:        collect.Version,
			CreateAt:       time.Time(collect.CreatedAt).UnixMilli(),
			UpdateAt:       time.Time(collect.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetEmojiCollectsByUuidsRes{
		Collects: collectItems,
	}, nil
}
