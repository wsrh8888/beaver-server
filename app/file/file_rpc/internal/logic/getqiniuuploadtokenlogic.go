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

package logic

import (
	"context"
	"time"

	"beaver/app/file/file_rpc/internal/svc"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetQiniuUploadTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQiniuUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQiniuUploadTokenLogic {
	return &GetQiniuUploadTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取七牛云上传token
func (l *GetQiniuUploadTokenLogic) GetQiniuUploadToken(in *file_rpc.GetQiniuUploadTokenReq) (*file_rpc.GetQiniuUploadTokenRes, error) {
	// 创建七牛云认证对象
	mac := qbox.NewMac(l.svcCtx.Config.Qiniu.AK, l.svcCtx.Config.Qiniu.SK)

	// 创建上传策略
	putPolicy := storage.PutPolicy{
		Scope:   l.svcCtx.Config.Qiniu.Bucket,                                                                 // 指定上传的目标资源空间
		Expires: uint64(time.Now().Add(time.Duration(l.svcCtx.Config.Qiniu.ExpireTime) * time.Second).Unix()), // 使用配置的过期时间
	}

	// 生成上传token
	uploadToken := putPolicy.UploadToken(mac)

	// 使用配置的过期时间（秒）
	expiresIn := l.svcCtx.Config.Qiniu.ExpireTime

	l.Logger.Infof("生成七牛云上传token成功，过期时间: %d秒", expiresIn)

	return &file_rpc.GetQiniuUploadTokenRes{
		UploadToken: uploadToken,
		ExpiresIn:   expiresIn,
	}, nil
}
