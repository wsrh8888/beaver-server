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
	"encoding/json"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	"github.com/google/uuid"
)

type CreateDigitalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateDigitalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDigitalLogic {
	return &CreateDigitalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("CreateDigitalLogic", ctx),
	}
}

func (l *CreateDigitalLogic) CreateDigital(in *open_rpc.CreateDigitalReq) (*open_rpc.CreateDigitalRes, error) {
	if in.Name == "" || in.OwnerId == "" {
		return nil, errors.New("名称和创建者不能为空")
	}
	if in.Runtime != open_models.DigitalRuntimeBuiltin && in.Runtime != open_models.DigitalRuntimeExternal {
		return nil, errors.New("runtime 无效")
	}
	if in.Source != open_models.DigitalSourceSystem && in.Source != open_models.DigitalSourceOpen {
		return nil, errors.New("source 无效")
	}
	visibility := in.Visibility
	if visibility == "" {
		visibility = open_models.DigitalVisibilityPrivate
	}
	if visibility != open_models.DigitalVisibilityPrivate &&
		visibility != open_models.DigitalVisibilityPartial &&
		visibility != open_models.DigitalVisibilityPublic {
		return nil, errors.New("visibility 无效")
	}

	createRes, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		NickName: in.Name,
		Abstract: in.Description,
		UserType: int32(user_models.UserTypeDigital),
		Source:   int32(user_models.SourceGroup),
	})
	if err != nil {
		return nil, errors.New("创建数字伙伴 IM 用户失败")
	}
	if in.Avatar != "" {
		_, _ = l.svcCtx.UserRpc.UserUpdateDisplay(l.ctx, &user_rpc.UserUpdateDisplayReq{
			UserId: createRes.UserID,
			Avatar: in.Avatar,
		})
	}

	partner := open_models.OpenDigital{
		DigitalID:  createRes.UserID,
		OwnerID:    in.OwnerId,
		Runtime:    in.Runtime,
		Visibility: visibility,
		Source:     in.Source,
		Status:     open_models.DigitalStatusEnabled,
	}
	res := &open_rpc.CreateDigitalRes{DigitalId: partner.DigitalID}
	if in.Runtime == open_models.DigitalRuntimeExternal {
		appKey, appSecret := newExternalCredentials()
		partner.AppKey = appKey
		partner.AppSecret = appSecret
		res.AppKey = appKey
		res.AppSecret = appSecret
	}
	if err := l.svcCtx.DB.Create(&partner).Error; err != nil {
		return nil, errors.New("创建数字伙伴失败")
	}
	if in.Runtime == open_models.DigitalRuntimeBuiltin {
		builtin, err := builtinConfigFrom(in.SystemPrompt, in.LlmModelId, in.Tools, in.McpServerIds)
		if err != nil {
			return nil, err
		}
		cfg := open_models.OpenDigitalConfig{DigitalID: partner.DigitalID, AppID: in.AppId, BuiltinConfig: builtin}
		if err := l.svcCtx.DB.Create(&cfg).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "保存数字伙伴配置失败", Data: map[string]any{"err": err.Error(), "digitalId": partner.DigitalID}})
			return nil, errors.New("保存数字伙伴配置失败")
		}
	}
	return res, nil
}

func newExternalCredentials() (appKey, appSecret string) {
	appKey = fmt.Sprintf("dig_%s", uuid.New().String()[:8])
	appSecret = uuid.New().String() + uuid.New().String()
	return appKey, appSecret
}

func builtinConfigFrom(prompt, modelID, tools, mcp string) (*open_models.BuiltinConfig, error) {
	cfg := &open_models.BuiltinConfig{SystemPrompt: prompt, LLMModelID: modelID}
	if tools != "" {
		if err := json.Unmarshal([]byte(tools), &cfg.Tools); err != nil {
			return nil, errors.New("tools 不是 JSON 数组")
		}
	}
	if mcp != "" {
		if err := json.Unmarshal([]byte(mcp), &cfg.McpServerIDs); err != nil {
			return nil, errors.New("mcp 不是 JSON 数组")
		}
	}
	return cfg, nil
}
