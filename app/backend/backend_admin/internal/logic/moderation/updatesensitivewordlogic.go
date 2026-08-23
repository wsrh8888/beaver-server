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
	"gorm.io/gorm"
)

type UpdateSensitiveWordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSensitiveWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSensitiveWordLogic {
	return &UpdateSensitiveWordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSensitiveWordLogic) UpdateSensitiveWord(req *types.UpdateSensitiveWordReq) (resp *types.UpdateSensitiveWordRes, err error) {
	var row backend_models.AdminSensitiveWord
	if err = l.svcCtx.DB.First(&row, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("敏感词不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if w := strings.TrimSpace(req.Word); w != "" {
		updates["word"] = w
	}
	if req.Category != "" {
		updates["category"] = strings.TrimSpace(req.Category)
	}
	if req.Level > 0 {
		updates["level"] = req.Level
	}
	if req.Remark != "" {
		updates["remark"] = strings.TrimSpace(req.Remark)
	}
	updates["is_active"] = req.IsActive

	if len(updates) == 0 {
		return &types.UpdateSensitiveWordRes{}, nil
	}

	if err = l.svcCtx.DB.Model(&row).Updates(updates).Error; err != nil {
		l.Errorf("更新敏感词失败: %v", err)
		return nil, errors.New("更新失败")
	}

	l.svcCtx.RecordOperation(req.UserID, "update_sensitive_word", "sensitive_word", row.Word, 0, "更新敏感词", "success", "")
	return &types.UpdateSensitiveWordRes{}, nil
}
