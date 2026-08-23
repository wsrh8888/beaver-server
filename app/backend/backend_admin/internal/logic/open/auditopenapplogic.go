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

package open

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type AuditOpenAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewAuditOpenAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditOpenAppLogic {
	return &AuditOpenAppLogic{logger: logger.New("audit_open_app"), ctx: ctx, svcCtx: svcCtx}
}

func (l *AuditOpenAppLogic) AuditOpenApp(req *types.AuditOpenAppReq) (resp *types.AuditOpenAppRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.AppID == "" {
		return nil, errors.New("应用ID不能为空")
	}
	if req.Status != 1 && req.Status != 2 {
		return nil, errors.New("无效的审核状态")
	}

	_, err = l.svcCtx.OpenRpc.UpdateOpenApps(l.ctx, &open_rpc.UpdateOpenAppsReq{
		AppIds:     []string{req.AppID},
		Action:     int32(req.Status),
		OperatorId: req.UserID,
		AuditRemark: req.AuditRemark,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("审核应用失败: %v", err)
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "开放应用审核成功",
		Data: map[string]interface{}{
			"appId":      req.AppID,
			"operatorId": req.UserID,
			"status":     req.Status,
		},
	})

	return &types.AuditOpenAppRes{}, nil
}
