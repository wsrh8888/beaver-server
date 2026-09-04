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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetQiniuUploadTokenLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取七牛云上传token
func NewGetQiniuUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQiniuUploadTokenLogic {
	return &GetQiniuUploadTokenLogic{
		logger: beaverlog.New("get_qiniu_upload_token", ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQiniuUploadTokenLogic) GetQiniuUploadToken(req *types.GetQiniuUploadTokenReq) (resp *types.GetQiniuUploadTokenRes, err error) {
	// 调用fileRpc服务获取七牛云上传token
	rpcResp, err := l.svcCtx.FileRpc.GetQiniuUploadToken(l.ctx, &file_rpc.GetQiniuUploadTokenReq{})
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "调用 fileRpc 获取七牛云 token 失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return nil, err
	}

	// 转换为HTTP响应格式
	resp = &types.GetQiniuUploadTokenRes{
		UploadToken: rpcResp.UploadToken,
		ExpiresIn:   rpcResp.ExpiresIn,
	}

	l.logger.Info(model.LogMsg{
		Text: "成功获取七牛云上传token",
		Data: map[string]interface{}{"expiresIn": resp.ExpiresIn},
	})
	return resp, nil
}
