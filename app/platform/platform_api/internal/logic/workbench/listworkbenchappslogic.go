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

package workbench

import (
	"context"
	"errors"
	"sort"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkbenchAppsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWorkbenchAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkbenchAppsLogic {
	return &ListWorkbenchAppsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func workbenchCategoryName(category int) string {
	switch category {
	case 1:
		return "办公"
	case 2:
		return "审批"
	case 3:
		return "效率"
	case 4:
		return "其他"
	default:
		return "默认"
	}
}

func (l *ListWorkbenchAppsLogic) ListWorkbenchApps(req *types.ListWorkbenchAppsReq) (*types.ListWorkbenchAppsRes, error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	rpcRes, err := l.svcCtx.PlatformRpc.ListEnabledWorkbenchApps(l.ctx, &platform_rpc.ListEnabledWorkbenchAppsReq{
		ClientScope: int32(req.ClientScope),
	})
	if err != nil {
		l.Errorf("获取工作台应用列表失败: %v", err)
		return nil, errors.New("获取工作台应用列表失败")
	}

	grouped := make(map[int][]types.ListWorkbenchAppsItem)
	for _, item := range rpcRes.List {
		entry := types.WorkbenchEntryConfig{}
		if item.EntryConfig != nil {
			entry = types.WorkbenchEntryConfig{
				Type:   int(item.EntryConfig.Type),
				PC:     item.EntryConfig.Pc,
				Mobile: item.EntryConfig.Mobile,
			}
		}
		category := int(item.Category)
		grouped[category] = append(grouped[category], types.ListWorkbenchAppsItem{
			WorkbenchAppID: item.WorkbenchAppId,
			Name:           item.Name,
			Description:    item.Description,
			Icon:           item.Icon,
			AppType:        int(item.AppType),
			ClientScope:    int(item.ClientScope),
			EntryConfig:    entry,
			Category:       category,
			Sort:           int(item.Sort),
			OpenMode:       int(item.OpenMode),
		})
	}

	categories := make([]int, 0, len(grouped))
	for category := range grouped {
		categories = append(categories, category)
	}
	sort.Ints(categories)

	groups := make([]types.ListWorkbenchAppsGroup, 0, len(categories))
	for _, category := range categories {
		groups = append(groups, types.ListWorkbenchAppsGroup{
			Category:     category,
			CategoryName: workbenchCategoryName(category),
			List:         grouped[category],
		})
	}

	return &types.ListWorkbenchAppsRes{Groups: groups}, nil
}
