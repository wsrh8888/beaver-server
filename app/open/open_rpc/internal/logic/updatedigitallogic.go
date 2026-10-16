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
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	beaverlog "beaver/utils/beaverlog"
	"gorm.io/gorm"
)

type UpdateDigitalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateDigitalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDigitalLogic {
	return &UpdateDigitalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("UpdateDigitalLogic", ctx),
	}
}

func (l *UpdateDigitalLogic) UpdateDigital(in *open_rpc.UpdateDigitalReq) (*open_rpc.UpdateDigitalRes, error) {
	if in.DigitalId == "" {
		return nil, errors.New("digitalId 不能为空")
	}
	if in.Status != int32(open_models.DigitalStatusDisabled) && in.Status != int32(open_models.DigitalStatusEnabled) {
		return nil, errors.New("status 无效")
	}
	if in.Visibility != "" &&
		in.Visibility != open_models.DigitalVisibilityPrivate &&
		in.Visibility != open_models.DigitalVisibilityPartial &&
		in.Visibility != open_models.DigitalVisibilityPublic {
		return nil, errors.New("visibility 无效")
	}

	var partner open_models.OpenDigital
	if err := l.svcCtx.DB.Where("digital_id = ?", in.DigitalId).First(&partner).Error; err != nil {
		return nil, errors.New("数字伙伴不存在")
	}

	if in.Name != "" {
		if _, err := l.svcCtx.UserRpc.UserUpdateDisplay(l.ctx, &user_rpc.UserUpdateDisplayReq{
			UserId:   in.DigitalId,
			NickName: in.Name,
		}); err != nil {
			return nil, errors.New("更新名称失败")
		}
	}

	updates := map[string]any{"status": in.Status}
	if in.Visibility != "" {
		updates["visibility"] = in.Visibility
	}
	if err := l.svcCtx.DB.Model(&partner).Updates(updates).Error; err != nil {
		return nil, errors.New("更新数字伙伴失败")
	}

	if partner.Runtime == open_models.DigitalRuntimeBuiltin && (in.SystemPrompt != "" || in.LlmModelId != "" || in.Tools != "" || in.McpServerIds != "") {
		builtin, err := builtinConfigFrom(in.SystemPrompt, in.LlmModelId, in.Tools, in.McpServerIds)
		if err != nil {
			return nil, err
		}
		var cfg open_models.OpenDigitalConfig
		err = l.svcCtx.DB.Where("digital_id = ?", in.DigitalId).First(&cfg).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cfg = open_models.OpenDigitalConfig{DigitalID: in.DigitalId, AppID: in.AppId, BuiltinConfig: builtin}
			if err := l.svcCtx.DB.Create(&cfg).Error; err != nil {
				return nil, errors.New("保存数字伙伴配置失败")
			}
		} else if err != nil {
			return nil, errors.New("读取编排配置失败")
		} else if err := l.svcCtx.DB.Model(&cfg).Updates(map[string]any{
			"app_id":         in.AppId,
			"builtin_config": builtin,
		}).Error; err != nil {
			return nil, errors.New("保存数字伙伴配置失败")
		}
		return &open_rpc.UpdateDigitalRes{}, nil
	}

	if in.AppId != "" && partner.Runtime == open_models.DigitalRuntimeBuiltin {
		var cfg open_models.OpenDigitalConfig
		err := l.svcCtx.DB.Where("digital_id = ?", in.DigitalId).First(&cfg).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cfg = open_models.OpenDigitalConfig{DigitalID: in.DigitalId, AppID: in.AppId}
			if err := l.svcCtx.DB.Create(&cfg).Error; err != nil {
				return nil, errors.New("绑定 AI 应用失败")
			}
		} else if err != nil {
			return nil, errors.New("读取编排配置失败")
		} else if err := l.svcCtx.DB.Model(&cfg).Update("app_id", in.AppId).Error; err != nil {
			return nil, errors.New("绑定 AI 应用失败")
		}
	}
	return &open_rpc.UpdateDigitalRes{}, nil
}
