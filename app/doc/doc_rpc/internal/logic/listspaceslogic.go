package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListSpacesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListSpacesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSpacesLogic {
	return &ListSpacesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListSpacesLogic", ctx),
	}
}

func (l *ListSpacesLogic) ListSpaces(in *doc_rpc.ListSpacesReq) (*doc_rpc.ListSpacesRes, error) {
	userID := trim(in.UserId)
	if userID == "" {
		// 用户隔离：无主体不返回任何空间。
		return &doc_rpc.ListSpacesRes{Total: 0}, nil
	}

	// 被显式授权（share/协作）的空间 ID 集合。
	var permitted []string
	l.svcCtx.DB.Model(&doc_models.DocPermission{}).
		Where("resource_type = ? AND principal_type = ? AND principal_id = ?", "space", "user", userID).
		Pluck("resource_id", &permitted)

	db := l.svcCtx.DB.Model(&doc_models.DocSpace{}).
		Where("deleted_at = ?", "1970-01-01 00:00:00").
		Where("owner_user_id = ? OR space_id IN ?", userID, append(permitted, ""))
	if in.Status > 0 {
		db = db.Where("status = ?", int(in.Status))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "统计空间失败", Data: err})
		return nil, err
	}

	page, pageSize := normalizePage(in.Page, in.PageSize)
	var rows []doc_models.DocSpace
	if err := db.Order("sort asc, id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "列出空间失败", Data: err})
		return nil, err
	}

	items := make([]*doc_rpc.SpaceItem, 0, len(rows))
	for i := range rows {
		items = append(items, toSpaceItem(&rows[i]))
	}
	return &doc_rpc.ListSpacesRes{Total: total, Spaces: items}, nil
}

// normalizePage 页码归一：page>=1，pageSize 默认 20、上限 200。
func normalizePage(page, size int32) (int, int) {
	p := int(page)
	if p < 1 {
		p = 1
	}
	s := int(size)
	if s <= 0 {
		s = 20
	}
	if s > 200 {
		s = 200
	}
	return p, s
}
