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

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAgentCredentialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAgentCredentialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentCredentialLogic {
	return &DeleteAgentCredentialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteAgentCredential 删通道。仍被任何模型引用时拒绝删除，避免出现指向空通道的孤儿模型。
func (l *DeleteAgentCredentialLogic) DeleteAgentCredential(in *platform_rpc.DeleteAgentCredentialReq) (*platform_rpc.DeleteAgentCredentialRes, error) {
	var refCount int64
	l.svcCtx.DB.Model(&agent_models.AgentModel{}).
		Where("credential_id = ?", in.CredentialId).Count(&refCount)
	if refCount > 0 {
		l.Logger.Errorf("接入通道仍被 %d 个模型引用，拒绝删除: credential_id=%s", refCount, in.CredentialId)
		return &platform_rpc.DeleteAgentCredentialRes{}, nil
	}
	l.svcCtx.DB.Where("credential_id = ?", in.CredentialId).Delete(&agent_models.AgentModelCredential{})
	return &platform_rpc.DeleteAgentCredentialRes{}, nil
}

