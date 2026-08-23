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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteSensitiveWordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSensitiveWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSensitiveWordLogic {
	return &DeleteSensitiveWordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSensitiveWordLogic) DeleteSensitiveWord(req *types.DeleteSensitiveWordReq) (resp *types.DeleteSensitiveWordRes, err error) {
	var row backend_models.AdminSensitiveWord
	if err = l.svcCtx.DB.First(&row, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("敏感词不存在")
		}
		return nil, err
	}

	if err = l.svcCtx.DB.Delete(&row).Error; err != nil {
		l.Errorf("删除敏感词失败: %v", err)
		return nil, err
	}

	l.svcCtx.RecordOperation(req.UserID, "delete_sensitive_word", "sensitive_word", row.Word, 0, "删除敏感词", "success", "")
	return &types.DeleteSensitiveWordRes{}, nil
}
