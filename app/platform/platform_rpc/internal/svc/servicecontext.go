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
	"beaver/app/platform/platform_rpc/internal/config"
	usercli "beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/database/platform/track"
	"beaver/database/platform/update"
	"beaver/utils/secret"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	UserRpc usercli.User
	// Secret 接入通道 API Key 加解密器，由 SecretKey 构建；启动时 mustCipher 校验，缺失即崩。
	Secret *secret.Cipher
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := coregorm.InitGorm(c.Mysql.DataSource)
	_ = update.InitUpdateApp(db)
	_ = track.InitBuckets(db)
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)
	return &ServiceContext{
		Config:  c,
		DB:      db,
		UserRpc: usercli.NewUser(zrpc.MustNewClient(c.UserRpc, rpcOpt)),
		Secret:  mustCipher(c.SecretKey),
	}
}

// mustCipher 用主密钥构造加解密器；任何错误都视为致命（未配密钥就无法解密通道 Key）。
func mustCipher(base64Key string) *secret.Cipher {
	cipher, err := secret.NewCipher(base64Key)
	if err != nil {
		panic(err)
	}
	return cipher
}
