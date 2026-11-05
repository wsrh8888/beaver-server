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

package digital

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	beaverlog "beaver/utils/beaverlog"
)

// goctl 会把 logx 写进 imports，占位避免未使用导入。
var _ = logx.WithContext

type GetDigitalConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取数字伙伴编排配置
func NewGetDigitalConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDigitalConfigLogic {
	return &GetDigitalConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetDigitalConfig", ctx),
	}
}

func (l *GetDigitalConfigLogic) GetDigitalConfig(req *types.GetDigitalConfigReq) (resp *types.GetDigitalConfigRes, err error) {
	if req.DigitalID == "" {
		return nil, errors.New("digitalId 不能为空")
	}
	if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
		return nil, err
	}

	// 查数字伙伴主表取 Runtime
	var digital open_models.OpenDigital
	if err := l.svcCtx.DB.Where("digital_id = ?", req.DigitalID).First(&digital).Error; err != nil {
		return nil, errors.New("数字伙伴不存在")
	}

	resp = &types.GetDigitalConfigRes{
		DigitalID: req.DigitalID,
		Runtime:   digital.Runtime,
	}

	// 查编排配置，按 Runtime 填对应字段
	var cfg open_models.OpenDigitalConfig
	if err := l.svcCtx.DB.Where("digital_id = ?", req.DigitalID).First(&cfg).Error; err == nil {
		switch digital.Runtime {
		case open_models.DigitalRuntimeBuiltin:
			if cfg.BuiltinConfig != nil {
				if b, mErr := json.Marshal(cfg.BuiltinConfig); mErr == nil {
					resp.BuiltinConfig = string(b)
				}
			}
		case open_models.DigitalRuntimeExternal:
			if cfg.ExternalConfig != nil {
				if b, mErr := json.Marshal(cfg.ExternalConfig); mErr == nil {
					resp.ExternalConfig = string(b)
				}
			}
		case open_models.DigitalRuntimeWorkflow:
			if cfg.WorkflowConfig != nil {
				if b, mErr := json.Marshal(cfg.WorkflowConfig); mErr == nil {
					resp.WorkflowConfig = string(b)
				}
			}
		}
	}

	return resp, nil
}
