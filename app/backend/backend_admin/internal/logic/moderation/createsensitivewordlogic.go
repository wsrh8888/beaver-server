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

package moderation

import (
	"context"
	"errors"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSensitiveWordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSensitiveWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSensitiveWordLogic {
	return &CreateSensitiveWordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSensitiveWordLogic) CreateSensitiveWord(req *types.CreateSensitiveWordReq) (resp *types.CreateSensitiveWordRes, err error) {
	word := strings.TrimSpace(req.Word)
	if word == "" {
		return nil, errors.New("敏感词不能为空")
	}
	level := req.Level
	if level <= 0 {
		level = 1
	}

	row := backend_models.AdminSensitiveWord{
		Word:     word,
		Category: strings.TrimSpace(req.Category),
		Level:    level,
		IsActive: true,
		Remark:   strings.TrimSpace(req.Remark),
	}
	if err = l.svcCtx.DB.Create(&row).Error; err != nil {
		l.Errorf("创建敏感词失败: %v", err)
		return nil, errors.New("创建失败，可能已存在相同词条")
	}

	l.svcCtx.RecordOperation(req.UserID, "create_sensitive_word", "sensitive_word", word, 0, "新增敏感词", "success", "")
	return &types.CreateSensitiveWordRes{ID: uint64(row.Id)}, nil
}
