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

package types

import "github.com/zeromicro/go-zero/core/logx"

type Config struct {
	Name      string
	Addr      string
	Etcd      string          // etcd 地址，如: 127.0.0.1:2379
	Auth      AuthConfig      // JWT 认证配置（网关直接解析 JWT，避免调用后端认证接口）
	Redis     RedisConfig     `json:",optional"` // Redis 配置（可选，用于验证 token 在 Redis 中的有效性）
	WhiteList []string        `json:",optional"` // 白名单路径，不需要认证的接口
	RateLimit RateLimitConfig `json:",optional"` // 限流配置
	Timeout   TimeoutConfig   `json:",optional"` // 超时配置
	Log       logx.LogConf    // 日志配置
}

type AuthConfig struct {
	AccessSecret string // JWT 密钥，必须与 backend_admin 的 AccessSecret 一致
	AccessExpire int    `json:",optional"` // token 过期时间（秒），可选
}

type RedisConfig struct {
	Addr     string `json:",optional"` // Redis 地址
	Password string `json:",optional"` // Redis 密码
	Db       int    `json:",optional"` // Redis 数据库编号
}

type RateLimitConfig struct {
	Enabled bool `json:",optional"` // 是否启用限流
	QPS     int  `json:",optional"` // 每秒请求数限制
}

type TimeoutConfig struct {
	Backend int `json:",optional"` // 后端请求超时（秒），默认 30
}
