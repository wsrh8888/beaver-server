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
	"beaver/app/call/call_rpc/internal/config"
	"beaver/core/coregorm"
	"beaver/core/coreredis"

	"beaver/app/chat/chat_rpc/chat"
	"beaver/app/user/user_rpc/user"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Redis   *redis.Client
	ChatRpc chat.Chat
	UserRpc user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	rdb := coreredis.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)

	return &ServiceContext{
		Config:  c,
		DB:      mysqlDb,
		Redis:   rdb,
		ChatRpc: chat.NewChat(zrpc.MustNewClient(c.ChatRpc)),
		UserRpc: user.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
