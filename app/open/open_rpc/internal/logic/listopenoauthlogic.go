/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package logic

import (
	"context"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type ListOpenOAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListOpenOAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOpenOAuthLogic {
	return &ListOpenOAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListOpenOAuthLogic", ctx),
	}
}

func (l *ListOpenOAuthLogic) ListOpenOAuth(in *open_rpc.ListOpenOAuthReq) (*open_rpc.ListOpenOAuthRes, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	query := l.svcCtx.DB.Model(&open_models.OpenOAuth{})
	if in.OwnerUserId != "" {
		query = query.Where("owner_id = ?", in.OwnerUserId)
	}
	if in.OauthId != "" {
		query = query.Where("oauth_id = ?", in.OauthId)
	}
	if in.Status != 0 {
		query = query.Where("status = ?", in.Status)
	}
	if in.Keyword != "" {
		like := "%" + in.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var apps []open_models.OpenOAuth
	offset := int((page - 1) * pageSize)
	if err := query.Order("id DESC").Offset(offset).Limit(int(pageSize)).Find(&apps).Error; err != nil {
		return nil, err
	}

	list := make([]*open_rpc.OpenOAuthItem, 0, len(apps))
	for _, app := range apps {
		list = append(list, toOpenOAuthItem(app))
	}
	return &open_rpc.ListOpenOAuthRes{Total: total, List: list}, nil
}
