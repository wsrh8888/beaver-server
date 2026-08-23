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

package monitor

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/core/coreonline"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOnlineStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 在线用户统计
func NewGetOnlineStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOnlineStatsLogic {
	return &GetOnlineStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOnlineStatsLogic) GetOnlineStats(req *types.GetOnlineStatsReq) (resp *types.GetOnlineStatsRes, err error) {
	online, err := coreonline.List(l.svcCtx.Redis)
	if err != nil {
		l.Errorf("获取在线用户统计失败: %v", err)
		return nil, err
	}

	resp = &types.GetOnlineStatsRes{
		UserCount: int64(len(online)),
	}

	for _, user := range online {
		for _, slot := range user.Slots {
			switch slot.Slot {
			case "desktop":
				resp.DesktopCount++
			case "mobile":
				resp.MobileCount++
			}
		}
	}

	return resp, nil
}
