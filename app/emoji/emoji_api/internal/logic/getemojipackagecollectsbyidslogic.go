package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageCollectsByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取用户收藏的表情包记录详情（同步用）
func NewGetEmojiPackageCollectsByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageCollectsByIdsLogic {
	return &GetEmojiPackageCollectsByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiPackageCollectsByIdsLogic) GetEmojiPackageCollectsByIds(req *types.GetEmojiPackageCollectsByIdsReq) (resp *types.GetEmojiPackageCollectsByIdsRes, err error) {
	if len(req.Ids) == 0 {
		return &types.GetEmojiPackageCollectsByIdsRes{
			Collects: make([]types.EmojiPackageCollectDetailItem, 0),
		}, nil
	}

	// 根据UUID列表查询收藏记录详情
	var collects []emoji_models.EmojiPackageCollect
	err = l.svcCtx.DB.Where("uuid IN ? AND user_id = ?", req.Ids, req.UserID).Find(&collects).Error
	if err != nil {
		l.Errorf("查询表情包收藏记录详情失败: uuids=%v, error=%v", req.Ids, err)
		return nil, err
	}

	l.Infof("批量查询表情包收藏记录详情: 请求%d个, 返回%d个", len(req.Ids), len(collects))

	// 转换为响应格式
	var collectItems []types.EmojiPackageCollectDetailItem
	for _, collect := range collects {
		collectItems = append(collectItems, types.EmojiPackageCollectDetailItem{
			UUID:      collect.UUID,
			UserID:    collect.UserID,
			PackageID: collect.PackageID,
			IsDeleted: collect.IsDeleted,
			Version:   collect.Version,
			CreateAt:  time.Time(collect.CreatedAt).UnixMilli(),
			UpdateAt:  time.Time(collect.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetEmojiPackageCollectsByIdsRes{
		Collects: collectItems,
	}, nil
}
