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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQiniuUploadTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取七牛云上传token
func NewGetQiniuUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQiniuUploadTokenLogic {
	return &GetQiniuUploadTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQiniuUploadTokenLogic) GetQiniuUploadToken(req *types.GetQiniuUploadTokenReq) (resp *types.GetQiniuUploadTokenRes, err error) {
	// 调用fileRpc服务获取七牛云上传token
	rpcResp, err := l.svcCtx.FileRpc.GetQiniuUploadToken(l.ctx, &file_rpc.GetQiniuUploadTokenReq{})
	if err != nil {
		l.Logger.Errorf("调用fileRpc获取七牛云token失败: %v", err)
		return nil, err
	}

	// 转换为HTTP响应格式
	resp = &types.GetQiniuUploadTokenRes{
		UploadToken: rpcResp.UploadToken,
		ExpiresIn:   rpcResp.ExpiresIn,
	}

	l.Logger.Infof("成功获取七牛云上传token，过期时间: %d秒", resp.ExpiresIn)
	return resp, nil
}
