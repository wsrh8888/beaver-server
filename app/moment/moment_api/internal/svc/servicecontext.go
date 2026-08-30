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
	"beaver/app/friend/friend_rpc/friend"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/moment/moment_api/internal/config"
	"beaver/app/notification/notification_rpc/notification"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	DB        *gorm.DB
	Redis     *redis.Client
	FriendRpc friend_rpc.FriendClient

	UserRpc   user_rpc.UserClient
	NotifyRpc notification.Notification
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	return &ServiceContext{
		Config:    c,
		DB:        mysqlDb,
		Redis:     client,
		UserRpc:   user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		FriendRpc: friend.NewFriend(zrpc.MustNewClient(c.FriendRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		NotifyRpc: notification.NewNotification(zrpc.MustNewClient(c.NotificationRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
