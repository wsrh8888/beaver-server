package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageContentsByRelationIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取表情包内容详情（通过relationIds，数据库同步）
func NewGetEmojiPackageContentsByRelationIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageContentsByRelationIdsLogic {
	return &GetEmojiPackageContentsByRelationIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiPackageContentsByRelationIdsLogic) GetEmojiPackageContentsByRelationIds(req *types.GetEmojiPackageContentsByRelationIdsReq) (resp *types.GetEmojiPackageContentsByRelationIdsRes, err error) {
	if len(req.RelationIds) == 0 {
		return &types.GetEmojiPackageContentsByRelationIdsRes{
			Contents: make([]types.EmojiPackageContentDetailItem, 0),
		}, nil
	}

	// 根据关联ID列表查询内容详情
	var contents []emoji_models.EmojiPackageEmoji
	err = l.svcCtx.DB.Where("relation_id IN ?", req.RelationIds).Find(&contents).Error
	if err != nil {
		l.Errorf("查询表情包内容详情失败: relationIds=%v, error=%v", req.RelationIds, err)
		return nil, err
	}

	l.Infof("批量查询表情包内容详情: 请求%d个关联ID, 返回%d条内容", len(req.RelationIds), len(contents))

	// 转换为响应格式
	var contentItems []types.EmojiPackageContentDetailItem
	for _, content := range contents {
		contentItems = append(contentItems, types.EmojiPackageContentDetailItem{
			RelationID: content.RelationID,
			PackageID:  content.PackageID,
			EmojiID:    content.EmojiID,
			SortOrder:  content.SortOrder,
			Version:    content.Version,
			CreatedAt:  time.Time(content.CreatedAt).UnixMilli(),
			UpdatedAt:  time.Time(content.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetEmojiPackageContentsByRelationIdsRes{
		Contents: contentItems,
	}, nil
}
