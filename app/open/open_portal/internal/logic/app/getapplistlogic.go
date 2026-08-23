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

package app

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取应用列表
func NewGetAppListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppListLogic {
	return &GetAppListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppListLogic) GetAppList(req *types.GetAppListReq) (resp *types.GetAppListRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}

	// 构建查询条件
	query := l.svcCtx.DB.Model(&open_models.OpenApp{}).Where("owner_user_id = ?", req.UserID)

	// 3. 如果指定了状态，添加状态过滤
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 4. 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	// 5. 分页查询
	var apps []open_models.OpenApp
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&apps).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	// 6. 转换为响应格式
	list := make([]types.AppInfo, 0, len(apps))
	for _, app := range apps {
		list = append(list, types.AppInfo{
			AppID:       app.AppID,
			Name:        app.Name,
			Description: app.Description,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt.Unix(),
		})
	}

	return &types.GetAppListRes{
		Total: total,
		List:  list,
	}, nil
}
