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
	aicli "beaver/remote/ai"

	"beaver/app/message/message_rpc/internal/config"

	agentcli "beaver/remote/agent"

	"beaver/app/relation/relation_rpc/client/friend"

	"beaver/app/relation/relation_rpc/client/group"

	"beaver/app/notification/notification_rpc/notification"
	opencli "beaver/remote/open"

	"beaver/app/platform/platform_rpc/platform"

	"beaver/app/account/account_rpc/client/user"

	"beaver/common/zrpc_interceptor"

	"beaver/core/corepush"

	"beaver/commonconfig"
	"beaver/core/coregorm"

	"beaver/core/coreredis"

	"beaver/core/corerocketmq"

	"beaver/core/corewebhook"

	"beaver/utils/beaverlog"

	versionPkg "beaver/core/version"

	"github.com/go-redis/redis"

	"github.com/zeromicro/go-zero/zrpc"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	DB *gorm.DB

	Redis *redis.Client

	VersionGen *versionPkg.VersionGenerator

	UserRpc user.User

	FriendRpc friend.Friend

	GroupRpc group.Group

	NotificationRpc notification.Notification

	PlatformRpc platform.Platform

	RocketMQ *corerocketmq.Client

	WebhookSender *corewebhook.WebhookSender

	PushSender *corepush.PushSender
}

func NewServiceContext(c config.Config) *ServiceContext {

	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)

	client := coreredis.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)

	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)

	userRpc := user.NewUser(zrpc.MustNewClient(c.UserRpc, rpcOpt))

	friendRpc := friend.NewFriend(zrpc.MustNewClient(c.FriendRpc, rpcOpt))

	groupRpc := group.NewGroup(zrpc.MustNewClient(c.GroupRpc, rpcOpt))

	notificationRpc := notification.NewNotification(zrpc.MustNewClient(c.NotificationRpc, rpcOpt))

	platformRpc := platform.NewPlatform(zrpc.MustNewClient(c.PlatformRpc, rpcOpt))

	// 装配各跨仓对端的 S2S 出站客户端（域名与凭据来自本服务 yaml）：
	// open=beaver-open 独立仓，ai=beaver-ai 独立仓，agent=beaver-agent(py，静态 X-Internal-Token)。
	commonconfig.Load(c.RemoteProxy)
	opencli.Init(c.RemoteProxy.Open.ClientID, c.RemoteProxy.Open.ClientSecret)
	aicli.Init(c.RemoteProxy.Ai.ClientID, c.RemoteProxy.Ai.ClientSecret)
	agentcli.Init(c.RemoteProxy.Agent.ClientID, c.RemoteProxy.Agent.ClientSecret)

	mqClient := corerocketmq.InitRocketMQ(c.RocketMQ.Addr)

	beaverlog.SetRocketMQ(mqClient)

	return &ServiceContext{

		Config: c,

		Redis: client,

		DB: mysqlDb,

		VersionGen: versionGen,

		UserRpc: userRpc,

		FriendRpc: friendRpc,

		GroupRpc: groupRpc,

		NotificationRpc: notificationRpc,

		PlatformRpc: platformRpc,

		RocketMQ: mqClient,

		WebhookSender: corewebhook.NewWebhookSender(

			corewebhook.Config{Timeout: 10, RetryCount: 3},

			newOpenWebhookLogWriter(),
		),

		PushSender: corepush.NewPushSender(corepush.Config{

			Enabled: c.Push.Enabled,

			FCM: corepush.FCMConfig{

				Enabled: c.Push.FCM.Enabled,

				ProjectID: c.Push.FCM.ProjectID,

				CredentialsFile: c.Push.FCM.CredentialsFile,
			},

			APNs: corepush.APNsConfig{

				Enabled: c.Push.APNs.Enabled,

				KeyFile: c.Push.APNs.KeyFile,

				KeyID: c.Push.APNs.KeyID,

				TeamID: c.Push.APNs.TeamID,

				BundleID: c.Push.APNs.BundleID,

				Production: c.Push.APNs.Production,
			},
		}),
	}

}
