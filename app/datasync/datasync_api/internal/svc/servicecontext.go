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
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"
	"beaver/app/datasync/datasync_api/internal/config"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/notification/notification_rpc/types/notification_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	"time"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	DB              *gorm.DB
	Redis           *redis.Client
	FriendRpc       friend_rpc.FriendClient
	GroupRpc        group_rpc.GroupClient
	CircleRpc       circle_rpc.CircleClient
	UserRpc         user_rpc.UserClient
	ChatRpc         chat_rpc.ChatClient
	EmojiRpc        emoji_rpc.EmojiClient
	NotificationRpc notification_rpc.NotificationClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	return &ServiceContext{
		Config:          c,
		DB:              mysqlDb,
		Redis:           client,
		FriendRpc:       friend_rpc.NewFriendClient(zrpc.MustNewClient(c.FriendRpc, zrpc.WithTimeout(time.Duration(c.FriendRpc.Timeout)*time.Millisecond)).Conn()),
		GroupRpc:        group_rpc.NewGroupClient(zrpc.MustNewClient(c.GroupRpc, zrpc.WithTimeout(time.Duration(c.GroupRpc.Timeout)*time.Millisecond)).Conn()),
		CircleRpc:       circle_rpc.NewCircleClient(zrpc.MustNewClient(c.CircleRpc, zrpc.WithTimeout(time.Duration(c.CircleRpc.Timeout)*time.Millisecond)).Conn()),
		UserRpc:         user_rpc.NewUserClient(zrpc.MustNewClient(c.UserRpc, zrpc.WithTimeout(time.Duration(c.UserRpc.Timeout)*time.Millisecond)).Conn()),
		ChatRpc:         chat_rpc.NewChatClient(zrpc.MustNewClient(c.ChatRpc, zrpc.WithTimeout(time.Duration(c.ChatRpc.Timeout)*time.Millisecond)).Conn()),
		EmojiRpc:        emoji_rpc.NewEmojiClient(zrpc.MustNewClient(c.EmojiRpc, zrpc.WithTimeout(time.Duration(c.EmojiRpc.Timeout)*time.Millisecond)).Conn()),
		NotificationRpc: notification_rpc.NewNotificationClient(zrpc.MustNewClient(c.NotificationRpc, zrpc.WithTimeout(time.Duration(c.NotificationRpc.Timeout)*time.Millisecond)).Conn()),
	}
}
