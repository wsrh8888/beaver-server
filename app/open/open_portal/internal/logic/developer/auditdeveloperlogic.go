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

package developer

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type AuditDeveloperLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewAuditDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditDeveloperLogic {
	return &AuditDeveloperLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("audit_developer", ctx)}
}

func (l *AuditDeveloperLogic) AuditDeveloper(req *types.AuditDeveloperReq) (resp *types.AuditDeveloperRes, err error) {
	userID, ok := l.ctx.Value("userId").(string)
	if !ok || userID == "" {
		return nil, errors.New("未登录")
	}
	if req.Status != 1 && req.Status != 2 {
		return nil, fmt.Errorf("无效的状态值")
	}

	_, err = l.svcCtx.OpenRpc.AuditDeveloper(l.ctx, &open_rpc.AuditDeveloperReq{
		Id:          uint64(req.ID),
		Status:      int32(req.Status),
		AuditBy:     userID,
		AuditRemark: req.AuditRemark,
	})
	if err != nil {
		return nil, err
	}

	l.logger.Info(model.LogMsg{Text: "开发者申请已审核", Data: map[string]interface{}{"id": req.ID, "status": req.Status}})
	return &types.AuditDeveloperRes{}, nil
}
