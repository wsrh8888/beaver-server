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
	"beaver/app/group/group_rpc/group"
	"beaver/app/open/open_portal/internal/config"
	"beaver/app/open/open_portal/internal/middleware"
	"beaver/app/open/open_rpc/open"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config                     config.Config
	DB                         *gorm.DB
	UserRpc                    user.User
	AuthRpc                    auth.Auth
	GroupRpc                   group.Group
	OpenRpc                    open.Open
	DeveloperAuthMiddleware    rest.Middleware
	RequireDeveloperMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := coregorm.InitGorm(c.Mysql.DataSource)
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)
	openRpc := open.NewOpen(zrpc.MustNewClient(c.OpenRpc, rpcOpt))
	return &ServiceContext{
		Config:                     c,
		DB:                         db,
		UserRpc:                    user.NewUser(zrpc.MustNewClient(c.UserRpc, rpcOpt)),
		AuthRpc:                    auth.NewAuth(zrpc.MustNewClient(c.AuthRpc, rpcOpt)),
		GroupRpc:                   group.NewGroup(zrpc.MustNewClient(c.GroupRpc, rpcOpt)),
		OpenRpc:                    openRpc,
		DeveloperAuthMiddleware:    middleware.NewDeveloperAuthMiddleware(c.Auth.AccessSecret).Handle,
		RequireDeveloperMiddleware: middleware.NewRequireDeveloperMiddleware(openRpc).Handle,
	}
}
