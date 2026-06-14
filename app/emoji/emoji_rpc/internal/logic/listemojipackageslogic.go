package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEmojiPackagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListEmojiPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEmojiPackagesLogic {
	return &ListEmojiPackagesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListEmojiPackagesLogic) ListEmojiPackages(in *emoji_rpc.ListEmojiPackagesReq) (*emoji_rpc.ListEmojiPackagesRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&emoji_models.EmojiPackage{})
	if in.PackageId != "" {
		db = db.Where("package_id = ?", in.PackageId)
	}
	if in.UserId != "" {
		db = db.Where("user_id = ?", in.UserId)
	}
	if in.Type != "" {
		db = db.Where("type = ?", in.Type)
	}
	if in.Status != 0 {
		db = db.Where("status = ?", in.Status)
	}
	if in.Title != "" {
		db = db.Where("title LIKE ? OR description LIKE ?", "%"+in.Title+"%", "%"+in.Title+"%")
	}
	if in.StartTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", in.StartTime); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if in.EndTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", in.EndTime); err == nil {
			db = db.Where("created_at <= ?", t)
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []emoji_models.EmojiPackage
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, err
	}

	items := make([]*emoji_rpc.EmojiPackageItem, 0, len(list))
	for _, p := range list {
		items = append(items, &emoji_rpc.EmojiPackageItem{
			PackageId:   p.PackageID,
			Title:       p.Title,
			CoverFile:   p.CoverFile,
			UserId:      p.UserID,
			Description: p.Description,
			Type:        p.Type,
			Status:      int32(p.Status),
			CreatedAt:   time.Time(p.CreatedAt).Format(time.RFC3339),
			UpdatedAt:   time.Time(p.UpdatedAt).Format(time.RFC3339),
		})
	}
	return &emoji_rpc.ListEmojiPackagesRes{Total: total, List: items}, nil
}
