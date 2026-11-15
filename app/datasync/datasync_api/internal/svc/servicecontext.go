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
	agentcli "beaver/remote/agent"
	"beaver/app/account/account_rpc/client/auth"
	assettypes "beaver/app/asset/asset_rpc/types"
	platform_rpc "beaver/app/platform/platform_rpc/types/platform_rpc"
	circle "beaver/remote/circle"
	"beaver/app/message/message_rpc/types/chat_rpc"
	"beaver/app/datasync/datasync_api/internal/config"
	"beaver/app/relation/relation_rpc/types"
	"beaver/app/notification/notification_rpc/types/notification_rpc"
	accttypes "beaver/app/account/account_rpc/types"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	"beaver/middleware/jwt"
	"time"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	DB              *gorm.DB
	Redis           *redis.Client
	FriendRpc       types.FriendClient
	GroupRpc        types.GroupClient
	UserRpc         accttypes.UserClient
	ChatRpc         chat_rpc.ChatClient
	EmojiRpc        assettypes.EmojiClient
	NotificationRpc notification_rpc.NotificationClient
	AuthRpc         auth.Auth
	FileRpc         assettypes.FileClient
	PlatformRpc     platform_rpc.PlatformClient
	// InternalToken 内部编排面（/internal/*）的 S2S Bearer 验签闸门（受众 server:internal）。
	InternalToken rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)

	// 构造期 fail-closed：S2S 密钥没配 / 注册表缺项，内部面整体不可用。
	if c.S2S.Secret == "" {
		panic("datasync_api S2S.Secret 未配置")
	}
	if len(c.S2S.Clients) == 0 {
		panic("datasync_api S2S.Clients 未配置")
	}

	// 装配跨仓对端的 S2S 出站客户端：circle=圈子服务（client_credentials 换票），
	// agent=beaver-agent(py，静态 X-Internal-Token)。域名与凭据来自本服务 yaml。
	circle.Init(c.Circle.ClientID, c.Circle.ClientSecret, c.Circle.BaseURL)
	agentcli.Init(c.Agent.BaseURL, c.Agent.InternalToken)

	return &ServiceContext{
		Config:          c,
		DB:              mysqlDb,
		Redis:           client,
		FriendRpc:       types.NewFriendClient(zrpc.MustNewClient(c.FriendRpc, zrpc.WithTimeout(time.Duration(c.FriendRpc.Timeout)*time.Millisecond)).Conn()),
		GroupRpc:        types.NewGroupClient(zrpc.MustNewClient(c.GroupRpc, zrpc.WithTimeout(time.Duration(c.GroupRpc.Timeout)*time.Millisecond)).Conn()),
		UserRpc:         accttypes.NewUserClient(zrpc.MustNewClient(c.UserRpc, zrpc.WithTimeout(time.Duration(c.UserRpc.Timeout)*time.Millisecond)).Conn()),
		ChatRpc:         chat_rpc.NewChatClient(zrpc.MustNewClient(c.ChatRpc, zrpc.WithTimeout(time.Duration(c.ChatRpc.Timeout)*time.Millisecond)).Conn()),
		EmojiRpc:        assettypes.NewEmojiClient(zrpc.MustNewClient(c.EmojiRpc, zrpc.WithTimeout(time.Duration(c.EmojiRpc.Timeout)*time.Millisecond)).Conn()),
		NotificationRpc: notification_rpc.NewNotificationClient(zrpc.MustNewClient(c.NotificationRpc, zrpc.WithTimeout(time.Duration(c.NotificationRpc.Timeout)*time.Millisecond)).Conn()),
		AuthRpc:         auth.NewAuth(zrpc.MustNewClient(c.AuthRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		FileRpc:         assettypes.NewFileClient(zrpc.MustNewClient(c.FileRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)).Conn()),
		PlatformRpc:     platform_rpc.NewPlatformClient(zrpc.MustNewClient(c.PlatformRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)).Conn()),
		InternalToken:   jwt.NewS2S(c.S2S.Secret, "server:internal"),
	}
}
