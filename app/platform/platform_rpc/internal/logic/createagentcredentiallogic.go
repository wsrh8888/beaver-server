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

type CreateAgentCredentialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentCredentialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentCredentialLogic {
	return &CreateAgentCredentialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateAgentCredential 新建接入通道。API Key 明文入参 → AES-GCM 加密落库 → 只回 credential_id，
// 密钥与密文都不回显。scope=user 必须带 user_id（个人通道归属）；缺省按 platform 官方通道。
func (l *CreateAgentCredentialLogic) CreateAgentCredential(in *platform_rpc.CreateAgentCredentialReq) (*platform_rpc.CreateAgentCredentialRes, error) {
	scope := in.Scope
	if scope == "" {
		scope = "platform"
	}
	if scope == "user" && in.UserId == "" {
		return &platform_rpc.CreateAgentCredentialRes{}, nil
	}
	if in.BaseUrl == "" || in.Name == "" {
		return &platform_rpc.CreateAgentCredentialRes{}, nil
	}
	enc, err := l.svcCtx.Secret.Encrypt(in.ApiKey)
	if err != nil {
		l.Logger.Errorf("接入通道密钥加密失败: %v", err)
		return &platform_rpc.CreateAgentCredentialRes{}, nil
	}
	var userID *string
	if in.UserId != "" {
		uid := in.UserId
		userID = &uid
	}
	channelType := in.ChannelType
	if channelType == "" {
		channelType = "openai_compatible"
	}
	modelService := in.ModelService
	if modelService == "" {
		modelService = "openai"
	}
	status := int(in.Status)
	if status == 0 {
		status = 1
	}
	credentialID := "achan_" + randHex(16)
	row := agent_models.AgentModelCredential{
		CredentialID:    credentialID,
		Name:            in.Name,
		ModelService:    modelService,
		ChannelType:     channelType,
		BaseURL:         in.BaseUrl,
		APIKeyEncrypted: enc,
		KeyHint:         secret.KeyHint(in.ApiKey),
		Scope:           scope,
		UserID:          userID,
		Sort:            int(in.Sort),
		Status:          status,
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		return &platform_rpc.CreateAgentCredentialRes{}, nil
	}
	return &platform_rpc.CreateAgentCredentialRes{CredentialId: credentialID}, nil
}
