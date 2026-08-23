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

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveWebhookLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveWebhookLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveWebhookLogLogic {
	return &SaveWebhookLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveWebhookLogLogic) SaveWebhookLog(in *open_rpc.SaveWebhookLogReq) (*open_rpc.SaveWebhookLogRes, error) {
	log := open_models.OpenWebhookLog{
		ConfigID:  in.ConfigId,
		AppID:     in.AppId,
		EventType: in.EventType,
		Status:    int(in.Status),
	}
	if err := l.svcCtx.DB.Create(&log).Error; err != nil {
		l.Errorf("保存 Webhook 日志失败: %v", err)
		return nil, err
	}
	return &open_rpc.SaveWebhookLogRes{}, nil
}
