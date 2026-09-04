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

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetSyncGroupRequestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取所有需要更新的入群申请版本
func NewGetSyncGroupRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncGroupRequestsLogic {
	return &GetSyncGroupRequestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_sync_group_requests", ctx),
	}
}

func (l *GetSyncGroupRequestsLogic) GetSyncGroupRequests(req *types.GetSyncGroupRequestsReq) (resp *types.GetSyncGroupRequestsRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.logger.Error(model.LogMsg{Text: "用户ID为空"})
		return nil, errors.New("用户ID不能为空")
	}

	// 获取用户群组申请版本信息
	versionResp, err := l.svcCtx.GroupRpc.GetUserGroupRequestVersions(l.ctx, &group_rpc.GetUserGroupRequestVersionsReq{
		UserID: userId,
		Since:  req.Since,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取用户群组申请版本失败", Data: map[string]any{"userId": userId, "err": err.Error()}})
		return nil, err
	}

	// 转换为响应格式，确保返回空数组而不是null
	groupVersions := make([]types.GroupRequestsVersionItem, 0)
	if versionResp.Versions != nil {
		for _, version := range versionResp.Versions {
			groupVersions = append(groupVersions, types.GroupRequestsVersionItem{
				GroupID: version.GroupID,
				Version: version.Version,
			})
		}
	}

	return &types.GetSyncGroupRequestsRes{
		GroupVersions:   groupVersions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
