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
	"beaver/app/auth/auth_rpc/auth"
	"beaver/app/chat/chat_rpc/chat"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_rpc/group"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/open/open_api/internal/config"
	oauthmiddle "beaver/app/open/open_api/internal/middle/oauth"
	"beaver/app/open/open_rpc/open"
	"beaver/app/open/open_rpc/types/open_rpc"
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
	Config   config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	OAuth    *oauthmiddle.Qrcode
	UserRpc  user_rpc.UserClient
	AuthRpc  auth.Auth
	ChatRpc  chat_rpc.ChatClient
	GroupRpc group_rpc.GroupClient
	OpenRpc  open_rpc.OpenClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)

	return &ServiceContext{
		Config:   c,
		DB:       mysqlDb,
		Redis:    client,
		OAuth:    oauthmiddle.NewQrcode(mysqlDb),
		UserRpc:  user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		AuthRpc:  auth.NewAuth(zrpc.MustNewClient(c.AuthRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		ChatRpc:  chat.NewChat(zrpc.MustNewClient(c.ChatRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		GroupRpc: group.NewGroup(zrpc.MustNewClient(c.GroupRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		OpenRpc:  open.NewOpen(zrpc.MustNewClient(c.OpenRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
