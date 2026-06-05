package logic

import (
	"context"
	"time"

	"beaver/app/moment/moment_models"
	"beaver/app/moment/moment_rpc/internal/svc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMomentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMomentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMomentsLogic {
	return &ListMomentsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListMomentsLogic) ListMoments(in *moment_rpc.ListMomentsReq) (*moment_rpc.ListMomentsRes, error) {
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

	db := l.svcCtx.DB.Model(&moment_models.MomentModel{})
	if in.MomentId != "" {
		db = db.Where("moment_id = ?", in.MomentId)
	}
	if in.UserId != "" {
		db = db.Where("user_id = ?", in.UserId)
	}
	if in.Keywords != "" {
		db = db.Where("content LIKE ?", "%"+in.Keywords+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计动态失败: %v", err)
		return nil, err
	}

	var list []moment_models.MomentModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询动态列表失败: %v", err)
		return nil, err
	}

	items := make([]*moment_rpc.MomentItem, 0, len(list))
	for _, m := range list {
		files := make([]*moment_rpc.MomentFileItem, 0)
		if m.Files != nil {
			for _, f := range *m.Files {
				files = append(files, &moment_rpc.MomentFileItem{FileKey: f.FileKey})
			}
		}
		items = append(items, &moment_rpc.MomentItem{
			MomentId:  m.MomentID,
			UserId:    m.UserID,
			Content:   m.Content,
			Files:     files,
			IsDeleted: m.IsDeleted,
			CreatedAt: time.Time(m.CreatedAt).Format(time.RFC3339),
			UpdatedAt: time.Time(m.UpdatedAt).Format(time.RFC3339),
		})
	}
	return &moment_rpc.ListMomentsRes{Total: total, List: items}, nil
}
