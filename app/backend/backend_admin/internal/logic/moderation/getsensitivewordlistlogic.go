package moderation

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSensitiveWordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSensitiveWordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSensitiveWordListLogic {
	return &GetSensitiveWordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSensitiveWordListLogic) GetSensitiveWordList(req *types.GetSensitiveWordListReq) (resp *types.GetSensitiveWordListRes, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	db := l.svcCtx.DB.Model(&backend_models.AdminSensitiveWord{})
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("word LIKE ? OR category LIKE ?", like, like)
	}
	if req.IsActive {
		db = db.Where("is_active = ?", true)
	}

	var total int64
	if err = db.Count(&total).Error; err != nil {
		l.Errorf("统计敏感词失败: %v", err)
		return nil, err
	}

	var rows []backend_models.AdminSensitiveWord
	if err = db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		l.Errorf("查询敏感词失败: %v", err)
		return nil, err
	}

	list := make([]types.SensitiveWordInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, types.SensitiveWordInfo{
			ID:        uint64(row.ID),
			Word:      row.Word,
			Category:  row.Category,
			Level:     row.Level,
			IsActive:  row.IsActive,
			Remark:    row.Remark,
			CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.GetSensitiveWordListRes{List: list, Total: total}, nil
}
