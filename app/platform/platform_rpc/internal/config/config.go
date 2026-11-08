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

package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	UserRpc zrpc.RpcClientConf
	// DocRpc 文档权威域 gRPC：知识库按 doc_id 读取在线文档快照。
	DocRpc zrpc.RpcClientConf
	// SecretKey 接入通道 API Key 的 AES-256-GCM 主密钥（32 字节 base64）。
	// 缺失会导致服务启动即失败，避免运行到解密调用点才崩。
	SecretKey string
	// OpenSearch 存知识分段向量本体；ChunkIndex 缺省时用 beaver-knowledge-chunks。
	OpenSearch struct {
		Addr       string `json:",optional"`
		ChunkIndex string `json:",optional"`
	}
	// Index 知识索引 worker 调参。
	Index struct {
		Enabled         bool `json:",optional,default=true"`
		IntervalSeconds int  `json:",optional,default=5"`
		BatchSize       int  `json:",optional,default=16"`
		MaxRetry        int  `json:",optional,default=3"`
		// Dims embedding 输出维度，建索引 mapping 用；改维度须重建索引。
		Dims int `json:",optional,default=1536"`
		// ChunkSize / Overlap 缺省分段参数（字符数），集合 index_config 可覆盖。
		ChunkSize int `json:",optional,default=800"`
		Overlap   int `json:",optional,default=80"`
	}
}
