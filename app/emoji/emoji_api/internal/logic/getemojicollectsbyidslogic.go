package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiCollectsByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取用户收藏的表情记录详情（同步用）
func NewGetEmojiCollectsByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiCollectsByIdsLogic {
	return &GetEmojiCollectsByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiCollectsByIdsLogic) GetEmojiCollectsByIds(req *types.GetEmojiCollectsByIdsReq) (resp *types.GetEmojiCollectsByIdsRes, err error) {
	if len(req.Ids) == 0 {
		return &types.GetEmojiCollectsByIdsRes{
			Collects: make([]types.EmojiCollectDetailItem, 0),
		}, nil
	}

	// 根据UUID列表查询收藏记录详情
	var collects []emoji_models.EmojiCollectEmoji
	err = l.svcCtx.DB.Where("uuid IN ? AND user_id = ?", req.Ids, req.UserID).Find(&collects).Error
	if err != nil {
		l.Errorf("查询表情收藏记录详情失败: uuids=%v, error=%v", req.Ids, err)
		return nil, err
	}

	l.Infof("批量查询表情收藏记录详情: 请求%d个, 返回%d个", len(req.Ids), len(collects))

	// 转换为响应格式
	var collectItems []types.EmojiCollectDetailItem
	for _, collect := range collects {
		collectItems = append(collectItems, types.EmojiCollectDetailItem{
			UUID:      collect.UUID,
			UserID:    collect.UserID,
			EmojiID:   collect.EmojiID,
			IsDeleted: collect.IsDeleted,
			Version:   collect.Version,
			CreateAt:  time.Time(collect.CreatedAt).UnixMilli(),
			UpdateAt:  time.Time(collect.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetEmojiCollectsByIdsRes{
		Collects: collectItems,
	}, nil
}
