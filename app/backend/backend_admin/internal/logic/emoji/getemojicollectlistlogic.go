package logic

import (
	"context"
	"strconv"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_models"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiCollectListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户收藏的表情图片列表
func NewGetEmojiCollectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiCollectListLogic {
	return &GetEmojiCollectListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiCollectListLogic) GetEmojiCollectList(req *types.GetEmojiCollectListReq) (resp *types.GetEmojiCollectListRes, err error) {
	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("1 = 1")

	// 按用户ID筛选
	if req.UserID != "" {
		whereClause = whereClause.Where("user_id = ?", req.UserID)
	}

	// 按表情ID筛选
	if req.EmojiID != "" {
		if emojiID, err := strconv.ParseUint(req.EmojiID, 10, 32); err == nil {
			whereClause = whereClause.Where("emoji_id = ?", emojiID)
		}
	}

	// 时间范围筛选
	if req.StartTime != "" {
		if startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime); err == nil {
			whereClause = whereClause.Where("created_at >= ?", startTime)
		}
	}

	if req.EndTime != "" {
		if endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime); err == nil {
			whereClause = whereClause.Where("created_at <= ?", endTime)
		}
	}

	// 分页参数校验
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 分页查询
	collects, count, err := list_query.ListQuery(l.svcCtx.DB, emoji_models.EmojiCollectEmoji{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  page,
			Limit: pageSize,
			Sort:  "created_at desc",
		},
		Where: whereClause,
	})

	if err != nil {
		logx.Errorf("查询表情收藏列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.GetEmojiCollectListItem
	for _, collect := range collects {
		emojiTitle := ""
		emojiFileName := ""
		// 通过 EmojiID 查询 Emoji 信息
		var emoji emoji_models.Emoji
		if err := l.svcCtx.DB.Where("uuid = ?", collect.EmojiID).First(&emoji).Error; err == nil {
			emojiTitle = emoji.Title
			emojiFileName = emoji.FileKey
		}

		list = append(list, types.GetEmojiCollectListItem{
			UUID:         collect.UUID,
			UserID:       collect.UserID,
			EmojiUUID:    collect.EmojiID,
			EmojiTitle:   emojiTitle,
			EmojiFileKey: emojiFileName,
			CreateTime:   collect.CreatedAt.String(),
			UpdateTime:   collect.UpdatedAt.String(),
		})
	}

	return &types.GetEmojiCollectListRes{
		List:  list,
		Total: count,
	}, nil
}
