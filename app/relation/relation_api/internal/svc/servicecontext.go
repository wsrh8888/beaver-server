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
 *
 * beaver-server-header-v1
 */

package svc

import (
	"beaver/app/message/message_rpc/chat"
	"beaver/app/message/message_rpc/types/chat_rpc"
	"beaver/app/notification/notification_rpc/notification"
	opencli "beaver/remote/open"
	"beaver/app/relation/relation_api/internal/config"
	"beaver/app/relation/relation_rpc/client/friend"
	"beaver/app/relation/relation_rpc/client/group"
	"beaver/app/relation/relation_rpc/types"
	accttypes "beaver/app/account/account_rpc/types"
	"beaver/app/account/account_rpc/client/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	"beaver/core/corerocketmq"
	"beaver/utils/beaverlog"
	versionPkg "beaver/core/version"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	Redis       *redis.Client
	UserRpc     accttypes.UserClient
	FriendRpc   types.FriendClient
	GroupRpc    types.GroupClient
	ChatRpc     chat_rpc.ChatClient
	NotifyRpc   notification.Notification
	VersionGen  *versionPkg.VersionGenerator
	RocketMQ    *corerocketmq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	// 初始化 RocketMQ 客户端
	mqClient := corerocketmq.InitRocketMQ(c.RocketMQ.Addr)
	beaverlog.SetRocketMQ(mqClient)

	// 装配开放接口服务（beaver-open 独立仓）S2S 出站客户端：域名与凭据来自本服务 yaml。
	opencli.Init(c.Open.ClientID, c.Open.ClientSecret, c.Open.BaseURL)

	return &ServiceContext{
		Config:      c,
		DB:          mysqlDb,
		Redis:       client,
		UserRpc:     user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		FriendRpc:   friend.NewFriend(zrpc.MustNewClient(c.RelationRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		GroupRpc:    group.NewGroup(zrpc.MustNewClient(c.RelationRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		ChatRpc:     chat.NewChat(zrpc.MustNewClient(c.ChatRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		NotifyRpc:   notification.NewNotification(zrpc.MustNewClient(c.NotificationRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		VersionGen:  versionGen,
		RocketMQ:    mqClient,
	}
}
