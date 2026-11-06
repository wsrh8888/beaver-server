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
	"beaver/utils/secret"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAgentCredentialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAgentCredentialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentCredentialLogic {
	return &UpdateAgentCredentialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateAgentCredential 改通道。credential_id / scope / user_id 建后不可改。
// api_key 传空=不改；有值才重算密文与 hint（只进不出）。
func (l *UpdateAgentCredentialLogic) UpdateAgentCredential(in *platform_rpc.UpdateAgentCredentialReq) (*platform_rpc.UpdateAgentCredentialRes, error) {
	var row agent_models.AgentModelCredential
	if err := l.svcCtx.DB.Where("credential_id = ?", in.CredentialId).First(&row).Error; err != nil {
		return &platform_rpc.UpdateAgentCredentialRes{}, nil
	}
	updates := map[string]any{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.ModelService != "" {
		updates["model_service"] = in.ModelService
	}
	if in.ChannelType != "" {
		updates["channel_type"] = in.ChannelType
	}
	if in.BaseUrl != "" {
		updates["base_url"] = in.BaseUrl
	}
	if in.ApiKey != "" {
		enc, err := l.svcCtx.Secret.Encrypt(in.ApiKey)
		if err != nil {
			l.Logger.Errorf("接入通道密钥加密失败: %v", err)
			return &platform_rpc.UpdateAgentCredentialRes{}, nil
		}
		updates["api_key_encrypted"] = enc
		updates["key_hint"] = secret.KeyHint(in.ApiKey)
	}
	if in.Sort != 0 {
		updates["sort"] = int(in.Sort)
	}
	if in.Status != 0 {
		updates["status"] = int(in.Status)
	}
	if len(updates) > 0 {
		if err := l.svcCtx.DB.Model(&agent_models.AgentModelCredential{}).
			Where("credential_id = ?", in.CredentialId).Updates(updates).Error; err != nil {
			return &platform_rpc.UpdateAgentCredentialRes{}, nil
		}
	}
	return &platform_rpc.UpdateAgentCredentialRes{}, nil
}
