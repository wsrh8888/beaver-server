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

type GetEmojiPackageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取表情包集合列表
func NewGetEmojiPackageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageListLogic {
	return &GetEmojiPackageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiPackageListLogic) GetEmojiPackageList(req *types.GetEmojiPackageListReq) (resp *types.GetEmojiPackageListRes, err error) {
	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("1 = 1")

	// 按创建者ID筛选
	if req.UserID != "" {
		whereClause = whereClause.Where("user_id = ?", req.UserID)
	}

	// 按类型筛选
	if req.Type != "" {
		whereClause = whereClause.Where("type = ?", req.Type)
	}

	// 按状态筛选
	if req.Status != 0 {
		whereClause = whereClause.Where("status = ?", req.Status)
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

	// 分页查询
	packages, count, err := list_query.ListQuery(l.svcCtx.DB, emoji_models.EmojiPackage{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.PageSize,
			Key:   req.Title,
			Sort:  "created_at desc",
		},
		Where: whereClause,
		Likes: []string{"title", "description"},
	})

	if err != nil {
		logx.Errorf("查询表情包列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.EmojiPackageInfo
	for _, pkg := range packages {
		list = append(list, types.EmojiPackageInfo{
			Id:          strconv.Itoa(int(pkg.Id)),
			Title:       pkg.Title,
			CoverFile:   pkg.CoverFile,
			UserID:      pkg.UserID,
			Description: pkg.Description,
			Type:        pkg.Type,
			Status:      pkg.Status,
			CreateTime:  pkg.CreatedAt.String(),
			UpdateTime:  pkg.UpdatedAt.String(),
		})
	}

	return &types.GetEmojiPackageListRes{
		List:  list,
		Total: count,
	}, nil
}
