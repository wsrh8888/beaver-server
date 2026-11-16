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

package config

import (
	"beaver/utils/s2s"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Etcd  string
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	RocketMQ struct {
		Addr string
	}
	Domain      string // 项目对外访问域名（用于生成本地文件完整URL）
	FileMaxSize map[string]float64
	WhiteList   []string
	BlackList   []string
	UserRpc     zrpc.RpcClientConf
	// AssetRpc：file + emoji 同进程（etcd assetrpc.rpc），队间能力面薄转发到它。
	AssetRpc zrpc.RpcClientConf
	// StorageType 通用上传接口(/beaver/api/file/file/v1/upload)使用的默认存储后端：local|qiniu|minio
	StorageType string
	Local       struct {
		UploadDir   string // 本地文件上传目录
		ProjectName string // 项目名称，用于文件路径前缀（为空则使用根目录）
	}
	Qiniu struct {
		ProjectName string // 项目名称，用于文件路径前缀（为空则使用根目录）
		AK          string
		SK          string
		Bucket      string
		Domain      string // 七牛云文件访问域名
		ExpireTime  int64
	}
	Minio struct {
		Endpoint    string // MinIO服务地址，如 127.0.0.1:9000
		AccessKey   string // MinIO访问密钥
		SecretKey   string // MinIO秘密密钥
		Bucket      string // MinIO存储桶名称
		UseSSL      bool   // 是否启用SSL
		ProjectName string // 项目名称，用于文件路径前缀（为空则使用根目录）
	}
	// Auth 用户 JWT 签名密钥，本服务独立签发并校验用户票：
	// 本服务独立验用户票，不信任任何外部直接塞进来的 Beaver-User-Id。
	Auth struct {
		AccessSecret string
		AccessExpire int
	}
	// S2S 入站机器票：队间调用与本服务共用同一套路由，凭据受众统一 core:api。
	S2S struct {
		Secret  string                     `json:",optional"`
		Clients map[string]s2s.ClientEntry `json:",optional"`
	}
}
