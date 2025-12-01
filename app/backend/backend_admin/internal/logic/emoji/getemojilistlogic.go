package logic

import (
	"context"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_models"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取表情图片列表
func NewGetEmojiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiListLogic {
	return &GetEmojiListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiListLogic) GetEmojiList(req *types.GetEmojiListReq) (resp *types.GetEmojiListRes, err error) {
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

	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("1 = 1")

	// 按创建者ID筛选
	// 作者ID筛选暂时移除

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

	// 分页查询
	emojis, count, err := list_query.ListQuery(l.svcCtx.DB, emoji_models.Emoji{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  page,
			Limit: pageSize,
			Key:   req.Title,
			Sort:  "created_at desc",
		},
		Where: whereClause,
		Likes: []string{"title"},
	})

	if err != nil {
		logx.Errorf("查询表情列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.GetEmojiListItem
	for _, emoji := range emojis {
		list = append(list, types.GetEmojiListItem{
			UUID:       emoji.UUID,
			FileKey:    emoji.FileKey,
			Title:      emoji.Title,
			AuthorID:   "", // 暂时为空，后续可从其他途径获取
			CreateTime: emoji.CreatedAt.String(),
			UpdateTime: emoji.UpdatedAt.String(),
		})
	}

	return &types.GetEmojiListRes{
		List:  list,
		Total: count,
	}, nil
}
