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

package overview

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/core/coreonline"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDashboardOverviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDashboardOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDashboardOverviewLogic {
	return &GetDashboardOverviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func countRPC(l *GetDashboardOverviewLogic, name string, fn func() (int64, error)) int64 {
	total, err := fn()
	if err != nil {
		l.Errorf("获取%s统计失败: %v", name, err)
		return 0
	}
	return total
}

func (l *GetDashboardOverviewLogic) GetDashboardOverview(req *types.GetDashboardOverviewReq) (resp *types.GetDashboardOverviewRes, err error) {
	resp = &types.GetDashboardOverviewRes{}

	resp.UserTotal = countRPC(l, "用户", func() (int64, error) {
		res, err := l.svcCtx.UserRpc.ListUsers(l.ctx, &user_rpc.ListUsersReq{Page: 1, PageSize: 1})
		if err != nil {
			return 0, err
		}
		return res.Total, nil
	})

	resp.GroupTotal = countRPC(l, "群组", func() (int64, error) {
		res, err := l.svcCtx.GroupRpc.ListGroups(l.ctx, &group_rpc.ListGroupsReq{Page: 1, PageSize: 1})
		if err != nil {
			return 0, err
		}
		return res.Total, nil
	})

	resp.FriendTotal = countRPC(l, "好友", func() (int64, error) {
		res, err := l.svcCtx.FriendRpc.ListFriends(l.ctx, &friend_rpc.ListFriendsReq{Page: 1, PageSize: 1})
		if err != nil {
			return 0, err
		}
		return res.Total, nil
	})

	resp.MessageTotal = countRPC(l, "消息", func() (int64, error) {
		res, err := l.svcCtx.ChatRpc.ListChatMessages(l.ctx, &chat_rpc.ListChatMessagesReq{Page: 1, PageSize: 1})
		if err != nil {
			return 0, err
		}
		return res.Total, nil
	})

	resp.MomentTotal = countRPC(l, "动态", func() (int64, error) {
		res, err := l.svcCtx.MomentRpc.ListMoments(l.ctx, &moment_rpc.ListMomentsReq{Page: 1, PageSize: 1})
		if err != nil {
			return 0, err
		}
		return res.Total, nil
	})

	resp.BlockTotal = countRPC(l, "黑名单", func() (int64, error) {
		res, err := l.svcCtx.FriendRpc.ListFriendBlocks(l.ctx, &friend_rpc.ListFriendBlocksReq{Page: 1, PageSize: 1})
		if err != nil {
			return 0, err
		}
		return res.Total, nil
	})

	resp.PendingDeveloperCount = countRPC(l, "待审开发者", func() (int64, error) {
		res, err := l.svcCtx.OpenRpc.ListDevelopers(l.ctx, &open_rpc.ListDevelopersReq{Page: 1, PageSize: 1, Status: 0})
		if err != nil {
			return 0, err
		}
		return res.Total, nil
	})

	resp.PendingAppCount = countRPC(l, "待审应用", func() (int64, error) {
		res, err := l.svcCtx.OpenRpc.ListOpenApps(l.ctx, &open_rpc.ListOpenAppsReq{Page: 1, PageSize: 1, AuditStatus: 0})
		if err != nil {
			return 0, err
		}
		return res.Total, nil
	})

	resp.PendingFeedbackCount = countRPC(l, "待处理反馈", func() (int64, error) {
		res, err := l.svcCtx.PlatformRpc.ListFeedback(l.ctx, &platform_rpc.ListFeedbackReq{Page: 1, PageSize: 1, Status: 1})
		if err != nil {
			return 0, err
		}
		return res.Total, nil
	})

	resp.PendingReportCount = countRPC(l, "待处理举报", func() (int64, error) {
		res, err := l.svcCtx.PlatformRpc.ListContentReports(l.ctx, &platform_rpc.ListContentReportsReq{Page: 1, PageSize: 1, Status: 1})
		if err != nil {
			return 0, err
		}
		return res.Total, nil
	})

	var pendingCases int64
	if err := l.svcCtx.DB.Model(&backend_models.AdminModerationCase{}).Where("status = ?", backend_models.CaseStatusPending).Count(&pendingCases).Error; err != nil {
		l.Errorf("获取待处理工单统计失败: %v", err)
	} else {
		resp.PendingCaseCount = pendingCases
	}

	resp.OnlineUserCount = countRPC(l, "在线用户", func() (int64, error) {
		online, err := coreonline.List(l.svcCtx.Redis)
		if err != nil {
			return 0, err
		}
		return int64(len(online)), nil
	})

	return resp, nil
}
