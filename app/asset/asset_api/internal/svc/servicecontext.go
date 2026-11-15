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

package svc

import (
	"beaver/app/asset/asset_api/internal/config"
	accttypes "beaver/app/account/account_rpc/types"
	"beaver/app/account/account_rpc/client/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	"beaver/core/corerocketmq"
	versionPkg "beaver/core/version"
	"beaver/utils/beaverlog"

	"github.com/go-redis/redis"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	UserRpc     accttypes.UserClient
	DB          *gorm.DB
	Redis       *redis.Client
	VersionGen  *versionPkg.VersionGenerator
	RocketMQ    *corerocketmq.Client
	MinioClient *minio.Client // MinIO客户端，仅在配置了Endpoint时初始化，供上传与预览预签名使用
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)

	mqClient := corerocketmq.InitRocketMQ(c.RocketMQ.Addr)
	beaverlog.SetRocketMQ(mqClient)

	svc := &ServiceContext{
		Config:     c,
		DB:         mysqlDb,
		Redis:      client,
		VersionGen: versionPkg.NewVersionGenerator(client, mysqlDb),
		RocketMQ:   mqClient,
		UserRpc:    user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}

	// 初始化MinIO客户端（仅在配置了Endpoint时）
	if c.Minio.Endpoint != "" {
		mClient, err := minio.New(c.Minio.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(c.Minio.AccessKey, c.Minio.SecretKey, ""),
			Secure: c.Minio.UseSSL,
		})
		if err == nil {
			svc.MinioClient = mClient
		}
	}

	return svc
}
