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
	"beaver/app/account/account_rpc/client/user"
	accttypes "beaver/app/account/account_rpc/types"
	"beaver/app/message/message_api/internal/config"
	"beaver/app/message/message_api/internal/middleware"
	"beaver/app/message/message_rpc/chat"
	"beaver/app/relation/relation_rpc/client/friend"
	"beaver/app/relation/relation_rpc/client/group"
	"beaver/app/relation/relation_rpc/types"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	"beaver/core/corerocketmq"
	versionPkg "beaver/core/version"
	"beaver/utils/beaverlog"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	Redis      *redis.Client
	UserRpc    accttypes.UserClient
	ChatRpc    chat.Chat
	FriendRpc  types.FriendClient
	GroupRpc   types.GroupClient
	DB         *gorm.DB
	VersionGen *versionPkg.VersionGenerator
	RocketMQ   *corerocketmq.Client
	// AuthGate 业务面的统一身份闸门：用户 JWT 与队间机器票共用同一套路由。
	AuthGate rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	// 初始化 RocketMQ 客户端

	mqClient := corerocketmq.InitRocketMQ(c.RocketMQ.Addr)
	beaverlog.SetRocketMQ(mqClient)

	return &ServiceContext{
		DB:         mysqlDb,
		Redis:      client,
		Config:     c,
		UserRpc:    user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		ChatRpc:    chat.NewChat(zrpc.MustNewClient(c.ChatRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		FriendRpc:  friend.NewFriend(zrpc.MustNewClient(c.FriendRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		GroupRpc:   group.NewGroup(zrpc.MustNewClient(c.GroupRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		VersionGen: versionGen,
		RocketMQ:   mqClient,
		AuthGate:   middleware.NewAuthGateMiddleware(c).Handle,
	}
}
