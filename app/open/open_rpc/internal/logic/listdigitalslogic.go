/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * 中文：
 * 本文件为海狸 IM（Beaver IM）开源项目源代码。
 * 版权所有 © 2024-2026 Beaver IM Team，基于 MIT 协议授权。
 * 禁止删除、篡改或替换本文件头部版权与许可声明。
 * 使用与商业授权说明：https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * English:
 * This file is part of the Beaver IM open-source project.
 * Copyright (c) 2024-2026 Beaver IM Team. Licensed under the MIT License.
 * Do not remove, alter, or replace this copyright and license header.
 * Usage & commercial licensing: https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package logic

import (
	"context"
	"encoding/json"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type ListDigitalsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListDigitalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDigitalsLogic {
	return &ListDigitalsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListDigitalsLogic", ctx),
	}
}

func (l *ListDigitalsLogic) ListDigitals(in *open_rpc.ListDigitalsReq) (*open_rpc.ListDigitalsRes, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := l.svcCtx.DB.Model(&open_models.OpenDigital{})
	if in.Source != "" {
		query = query.Where("source = ?", in.Source)
	}
	if in.Status >= 0 {
		query = query.Where("status = ?", in.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var rows []open_models.OpenDigital
	if err := query.Order("id desc").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&rows).Error; err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.DigitalID)
	}

	names := map[string]*user_rpc.UserInfo{}
	if len(ids) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: ids}); err == nil && res != nil {
			names = res.UserInfo
		}
	}

	cfgByID := map[string]open_models.OpenDigitalConfig{}
	if len(ids) > 0 {
		var cfgs []open_models.OpenDigitalConfig
		_ = l.svcCtx.DB.Where("digital_id IN ?", ids).Find(&cfgs).Error
		for _, cfg := range cfgs {
			cfgByID[cfg.DigitalID] = cfg
		}
	}

	list := make([]*open_rpc.DigitalItem, 0, len(rows))
	for _, row := range rows {
		cfg := cfgByID[row.DigitalID]
		item := &open_rpc.DigitalItem{
			DigitalId:  row.DigitalID,
			OwnerId:    row.OwnerID,
			Runtime:    row.Runtime,
			Visibility: row.Visibility,
			Status:     int32(row.Status),
			Source:     row.Source,
			AppId:      cfg.AppID,
			CreatedAt:  row.CreatedAt.Unix(),
		}
		if cfg.BuiltinConfig != nil {
			item.SystemPrompt = cfg.BuiltinConfig.SystemPrompt
			item.LlmModelId = cfg.BuiltinConfig.LLMModelID
			if b, err := json.Marshal(cfg.BuiltinConfig.Tools); err == nil {
				item.Tools = string(b)
			}
			if b, err := json.Marshal(cfg.BuiltinConfig.McpServerIDs); err == nil {
				item.McpServerIds = string(b)
			}
		}
		if user, ok := names[row.DigitalID]; ok && user != nil {
			item.Name = user.NickName
			item.Description = user.Abstract
			item.Avatar = user.Avatar
		}
		list = append(list, item)
	}
	return &open_rpc.ListDigitalsRes{Total: total, List: list}, nil
}
