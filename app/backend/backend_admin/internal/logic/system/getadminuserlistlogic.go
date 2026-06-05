package system

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdminUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdminUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminUserListLogic {
	return &GetAdminUserListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAdminUserListLogic) GetAdminUserList(req *types.GetAdminUserListReq) (resp *types.GetAdminUserListRes, err error) {
	page, pageSize := req.Page, req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	db := l.svcCtx.DB.Model(&backend_models.AdminUser{})
	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		db = db.Where("nick_name LIKE ? OR phone LIKE ? OR user_id LIKE ?", kw, kw, kw)
	}
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	}

	var total int64
	if err = db.Count(&total).Error; err != nil {
		l.Errorf("统计管理员失败: %v", err)
		return nil, err
	}

	var rows []backend_models.AdminUser
	if err = db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		l.Errorf("查询管理员失败: %v", err)
		return nil, err
	}

	userIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.UserID)
	}

	authMap := make(map[string][]uint)
	nameMap := make(map[uint]string)
	if len(userIDs) > 0 {
		var authUsers []backend_models.AdminSystemAuthorityUser
		_ = l.svcCtx.DB.Where("user_id IN ?", userIDs).Find(&authUsers).Error
		authIDs := make([]uint, 0)
		for _, au := range authUsers {
			authMap[au.UserID] = append(authMap[au.UserID], au.AuthorityID)
			authIDs = append(authIDs, au.AuthorityID)
		}
		if len(authIDs) > 0 {
			var authorities []backend_models.AdminSystemAuthority
			_ = l.svcCtx.DB.Where("id IN ?", authIDs).Find(&authorities).Error
			for _, a := range authorities {
				nameMap[uint(a.Id)] = a.Name
			}
		}
	}

	list := make([]types.AdminUserInfo, 0, len(rows))
	for _, row := range rows {
		ids := authMap[row.UserID]
		names := make([]string, 0, len(ids))
		for _, id := range ids {
			if n, ok := nameMap[id]; ok {
				names = append(names, n)
			}
		}
		list = append(list, types.AdminUserInfo{
			Id:             uint(row.Id),
			UserID:         row.UserID,
			NickName:       row.NickName,
			Phone:          row.Phone,
			Status:         row.Status,
			LastLoginAt:    row.LastLoginAt,
			CreatedAt:      row.CreatedAt.String(),
			AuthorityIds:   ids,
			AuthorityNames: names,
		})
	}
	return &types.GetAdminUserListRes{List: list, Total: total}, nil
}
