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

package svc

import (
	"beaver/app/chat/chat_rpc/internal/config"
	"beaver/app/friend/friend_rpc/friend"
	"beaver/app/group/group_rpc/group"
	"beaver/app/notification/notification_rpc/notification"
	"beaver/app/open/open_rpc/open"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/corepush"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	"beaver/core/corerocketmq"
	"beaver/core/corewebhook"
	versionPkg "beaver/core/version"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	DB              *gorm.DB
	Redis           *redis.Client
	VersionGen      *versionPkg.VersionGenerator
	UserRpc         user.User
	FriendRpc       friend.Friend
	GroupRpc        group.Group
	OpenRpc         open.Open
	NotificationRpc notification.Notification
	RocketMQ        *corerocketmq.Client
	WebhookSender   *corewebhook.WebhookSender
	PushSender      *corepush.PushSender
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)
	userRpc := user.NewUser(zrpc.MustNewClient(c.UserRpc, rpcOpt))
	friendRpc := friend.NewFriend(zrpc.MustNewClient(c.FriendRpc, rpcOpt))
	groupRpc := group.NewGroup(zrpc.MustNewClient(c.GroupRpc, rpcOpt))
	openRpc := open.NewOpen(zrpc.MustNewClient(c.OpenRpc, rpcOpt))
	notificationRpc := notification.NewNotification(zrpc.MustNewClient(c.NotificationRpc, rpcOpt))
	mqClient := corerocketmq.InitRocketMQ(c.RocketMQ.Addr)

	return &ServiceContext{
		Config:          c,
		Redis:           client,
		DB:              mysqlDb,
		VersionGen:      versionGen,
		UserRpc:         userRpc,
		FriendRpc:       friendRpc,
		GroupRpc:        groupRpc,
		OpenRpc:         openRpc,
		NotificationRpc: notificationRpc,
		RocketMQ:        mqClient,
		WebhookSender: corewebhook.NewWebhookSender(
			corewebhook.Config{Timeout: 10, RetryCount: 3},
			newOpenRpcWebhookLogWriter(openRpc),
		),
		PushSender: corepush.NewPushSender(corepush.Config{
			Enabled: c.Push.Enabled,
			FCM: corepush.FCMConfig{
				Enabled:         c.Push.FCM.Enabled,
				ProjectID:       c.Push.FCM.ProjectID,
				CredentialsFile: c.Push.FCM.CredentialsFile,
			},
			APNs: corepush.APNsConfig{
				Enabled:    c.Push.APNs.Enabled,
				KeyFile:    c.Push.APNs.KeyFile,
				KeyID:      c.Push.APNs.KeyID,
				TeamID:     c.Push.APNs.TeamID,
				BundleID:   c.Push.APNs.BundleID,
				Production: c.Push.APNs.Production,
			},
		}),
	}
}
