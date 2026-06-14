package logic

import (
	"context"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListContentReportsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListContentReportsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListContentReportsLogic {
	return &ListContentReportsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListContentReportsLogic) ListContentReports(in *platform_rpc.ListContentReportsReq) (*platform_rpc.ListContentReportsRes, error) {
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

	db := l.svcCtx.DB.Model(&platform_models.ContentReportModel{})
	if in.Status > 0 {
		db = db.Where("status = ?", in.Status)
	}
	if in.TargetType > 0 {
		db = db.Where("target_type = ?", in.TargetType)
	}
	if in.TargetId != "" {
		db = db.Where("target_id = ?", in.TargetId)
	}
	if in.ReporterUserId != "" {
		db = db.Where("reporter_user_id = ?", in.ReporterUserId)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []platform_models.ContentReportModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, err
	}

	items := make([]*platform_rpc.ContentReportItem, 0, len(list))
	for _, r := range list {
		items = append(items, toContentReportItem(r))
	}
	return &platform_rpc.ListContentReportsRes{Total: total, List: items}, nil
}
