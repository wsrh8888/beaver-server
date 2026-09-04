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

type GetSyncGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取所有需要更新的群组信息版本
func NewGetSyncGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncGroupInfoLogic {
	return &GetSyncGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_sync_group_info", ctx),
	}
}

func (l *GetSyncGroupInfoLogic) GetSyncGroupInfo(req *types.GetSyncGroupInfoReq) (resp *types.GetSyncGroupInfoRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.logger.Error(model.LogMsg{Text: "用户ID为空"})
		return nil, errors.New("用户ID不能为空")
	}

	// 1. 获取用户加入的群组ID列表
	groupIDsResp, err := l.svcCtx.GroupRpc.GetUserGroupIDs(l.ctx, &group_rpc.GetUserGroupIDsReq{
		UserID: userId,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取用户群组ID列表失败", Data: map[string]any{"userId": userId, "err": err.Error()}})
		return nil, err
	}

	groupIDs := groupIDsResp.GroupIDs
	if len(groupIDs) == 0 {
		return &types.GetSyncGroupInfoRes{
			GroupVersions:   []types.GroupInfoVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	// 2. 获取变更的群组资料
	serverTimestamp := time.Now().UnixMilli()

	groupResp, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{
		GroupIDs: groupIDs,
		Since:    req.Since,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取变更的群组资料失败", Data: map[string]any{"userId": userId, "err": err.Error()}})
		return nil, err
	}

	// 3. 转换为响应格式，确保返回空数组而不是null
	groupVersions := make([]types.GroupInfoVersionItem, 0)
	if groupResp.Groups != nil {
		for _, group := range groupResp.Groups {
			groupVersions = append(groupVersions, types.GroupInfoVersionItem{
				GroupID: group.GroupID,
				Version: group.Version,
			})
		}
	}

	return &types.GetSyncGroupInfoRes{
		GroupVersions:   groupVersions,
		ServerTimestamp: serverTimestamp,
	}, nil
}
