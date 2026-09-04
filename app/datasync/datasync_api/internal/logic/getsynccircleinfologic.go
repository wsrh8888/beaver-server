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
	"errors"
	"time"

	"beaver/app/circle/circle_rpc/types/circle_rpc"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetSyncCircleInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取所有需要更新的圈子信息版本
func NewGetSyncCircleInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncCircleInfoLogic {
	return &GetSyncCircleInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_sync_circle_info", ctx),
	}
}

func (l *GetSyncCircleInfoLogic) GetSyncCircleInfo(req *types.GetSyncCircleInfoReq) (resp *types.GetSyncCircleInfoRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.logger.Error(model.LogMsg{Text: "用户ID为空"})
		return nil, errors.New("用户ID不能为空")
	}

	serverTimestamp := time.Now().UnixMilli()

	circleResp, err := l.svcCtx.CircleRpc.GetCircleVersions(l.ctx, &circle_rpc.GetCircleVersionsReq{
		UserId:  userId,
		Version: req.Since,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取变更的圈子资料失败", Data: map[string]any{"userId": userId, "err": err.Error()}})
		return nil, err
	}

	circleVersions := make([]types.CircleInfoVersionItem, 0)
	if circleResp.List != nil {
		for _, item := range circleResp.List {
			circleVersions = append(circleVersions, types.CircleInfoVersionItem{
				CircleID: item.CircleId,
				Version:  item.Version,
			})
		}
	}

	return &types.GetSyncCircleInfoRes{
		CircleVersions:  circleVersions,
		ServerTimestamp: serverTimestamp,
	}, nil
}
