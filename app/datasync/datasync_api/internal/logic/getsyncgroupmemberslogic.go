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
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncGroupMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的群成员版本
func NewGetSyncGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncGroupMembersLogic {
	return &GetSyncGroupMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncGroupMembersLogic) GetSyncGroupMembers(req *types.GetSyncGroupMembersReq) (resp *types.GetSyncGroupMembersRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	// 1. 获取用户加入的群组ID列表
	groupIDsResp, err := l.svcCtx.GroupRpc.GetUserGroupIDs(l.ctx, &group_rpc.GetUserGroupIDsReq{
		UserID: userId,
	})
	if err != nil {
		l.Errorf("获取用户群组ID列表失败: %v", err)
		return nil, err
	}

	groupIDs := groupIDsResp.GroupIDs
	if len(groupIDs) == 0 {
		return &types.GetSyncGroupMembersRes{
			GroupVersions:   []types.GroupMembersVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	// 2. 获取变更的群成员
	serverTimestamp := time.Now().UnixMilli()

	memberResp, err := l.svcCtx.GroupRpc.GetGroupMembersListByIds(l.ctx, &group_rpc.GetGroupMembersListByIdsReq{
		GroupIDs: groupIDs,
		Since:    req.Since,
	})
	if err != nil {
		l.Errorf("获取变更的群成员失败: %v", err)
		return nil, err
	}

	// 3. 合并版本信息，按群组聚合最新版本
	groupVersionsMap := make(map[string]int64)
	for _, member := range memberResp.Members {
		if currentVersion, exists := groupVersionsMap[member.GroupID]; !exists || member.Version > currentVersion {
			groupVersionsMap[member.GroupID] = member.Version
		}
	}

	// 4. 转换为响应格式，确保返回空数组而不是null
	groupVersions := make([]types.GroupMembersVersionItem, 0)
	if groupVersionsMap != nil {
		for groupID, version := range groupVersionsMap {
			groupVersions = append(groupVersions, types.GroupMembersVersionItem{
				GroupID: groupID,
				Version: version,
			})
		}
	}

	return &types.GetSyncGroupMembersRes{
		GroupVersions:   groupVersions,
		ServerTimestamp: serverTimestamp,
	}, nil
}
