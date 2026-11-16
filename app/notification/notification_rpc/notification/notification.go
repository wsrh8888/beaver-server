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

// 通知域聚合客户端：把 notification_rpc 拆出的 event/inbox/push 三个 service
// 组合成一个 Notification 接口，供仓内各服务复用，避免每个调用方分别持有多个 client。
// 这是手写维护的聚合层（非 goctl 生成物）。

package notification

import (
	"context"

	eventcli "beaver/app/notification/notification_rpc/client/event"
	inboxcli "beaver/app/notification/notification_rpc/client/inbox"
	pushcli "beaver/app/notification/notification_rpc/client/push"
	"beaver/app/notification/notification_rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	DisablePushTokenReq      = pb.DisablePushTokenReq
	DisablePushTokenRes      = pb.DisablePushTokenRes
	EventVersion             = pb.EventVersion
	GetEventVersionsReq      = pb.GetEventVersionsReq
	GetEventVersionsRes      = pb.GetEventVersionsRes
	GetInboxVersionsReq      = pb.GetInboxVersionsReq
	GetInboxVersionsRes      = pb.GetInboxVersionsRes
	GetReadCursorVersionsReq = pb.GetReadCursorVersionsReq
	GetReadCursorVersionsRes = pb.GetReadCursorVersionsRes
	InboxVersion             = pb.InboxVersion
	ListPushTokensReq        = pb.ListPushTokensReq
	ListPushTokensRes        = pb.ListPushTokensRes
	PushEventReq             = pb.PushEventReq
	PushEventRes             = pb.PushEventRes
	PushTokenInfo            = pb.PushTokenInfo
	ReadCursorVersion        = pb.ReadCursorVersion

	Notification interface {
		PushEvent(ctx context.Context, in *PushEventReq, opts ...grpc.CallOption) (*PushEventRes, error)
		GetEventVersions(ctx context.Context, in *GetEventVersionsReq, opts ...grpc.CallOption) (*GetEventVersionsRes, error)
		GetInboxVersions(ctx context.Context, in *GetInboxVersionsReq, opts ...grpc.CallOption) (*GetInboxVersionsRes, error)
		GetReadCursorVersions(ctx context.Context, in *GetReadCursorVersionsReq, opts ...grpc.CallOption) (*GetReadCursorVersionsRes, error)
		ListPushTokens(ctx context.Context, in *ListPushTokensReq, opts ...grpc.CallOption) (*ListPushTokensRes, error)
		DisablePushToken(ctx context.Context, in *DisablePushTokenReq, opts ...grpc.CallOption) (*DisablePushTokenRes, error)
	}

	defaultNotification struct {
		event eventcli.Event
		inbox inboxcli.Inbox
		push  pushcli.Push
	}
)

func NewNotification(cli zrpc.Client) Notification {
	return &defaultNotification{
		event: eventcli.NewEvent(cli),
		inbox: inboxcli.NewInbox(cli),
		push:  pushcli.NewPush(cli),
	}
}

func (m *defaultNotification) PushEvent(ctx context.Context, in *PushEventReq, opts ...grpc.CallOption) (*PushEventRes, error) {
	return m.event.PushEvent(ctx, in, opts...)
}

func (m *defaultNotification) GetEventVersions(ctx context.Context, in *GetEventVersionsReq, opts ...grpc.CallOption) (*GetEventVersionsRes, error) {
	return m.event.GetEventVersions(ctx, in, opts...)
}

func (m *defaultNotification) GetInboxVersions(ctx context.Context, in *GetInboxVersionsReq, opts ...grpc.CallOption) (*GetInboxVersionsRes, error) {
	return m.inbox.GetInboxVersions(ctx, in, opts...)
}

func (m *defaultNotification) GetReadCursorVersions(ctx context.Context, in *GetReadCursorVersionsReq, opts ...grpc.CallOption) (*GetReadCursorVersionsRes, error) {
	return m.inbox.GetReadCursorVersions(ctx, in, opts...)
}

func (m *defaultNotification) ListPushTokens(ctx context.Context, in *ListPushTokensReq, opts ...grpc.CallOption) (*ListPushTokensRes, error) {
	return m.push.ListPushTokens(ctx, in, opts...)
}

func (m *defaultNotification) DisablePushToken(ctx context.Context, in *DisablePushTokenReq, opts ...grpc.CallOption) (*DisablePushTokenRes, error) {
	return m.push.DisablePushToken(ctx, in, opts...)
}
