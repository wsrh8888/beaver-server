package system

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuthorityListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAuthorityListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorityListLogic {
	return &GetAuthorityListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAuthorityListLogic) GetAuthorityList(req *types.GetAuthorityListReq) (resp *types.GetAuthorityListRes, err error) {
	var rows []backend_models.AdminSystemAuthority
	if err = l.svcCtx.DB.Order("sort ASC, id ASC").Find(&rows).Error; err != nil {
		l.Errorf("查询角色列表失败: %v", err)
		return nil, err
	}

	list := make([]types.AuthorityInfo, 0, len(rows))
	for _, row := range rows {
		var menuCount int64
		_ = l.svcCtx.DB.Model(&backend_models.AdminSystemAuthorityMenu{}).
			Where("authority_id = ?", row.Id).Count(&menuCount).Error
		list = append(list, types.AuthorityInfo{
			Id:          uint(row.Id),
			Name:        row.Name,
			Description: row.Description,
			Status:      row.Status,
			Sort:        row.Sort,
			MenuCount:   menuCount,
		})
	}
	return &types.GetAuthorityListRes{List: list}, nil
}
