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
	"fmt"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"gorm.io/gorm"
)


type UpdateEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewUpdateEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmailLogic {
	return &UpdateEmailLogic{
		ctx:    ctx,
		logger: logger.New("update_email"),
		svcCtx: svcCtx,
	}
}

func (l *UpdateEmailLogic) UpdateEmail(req *types.UpdateEmailReq) (resp *types.UpdateEmailRes, err error) {
	// 获取要更新的用户信息
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	// 检查新邮箱是否与当前邮箱相同
	if user.Email == req.Email {
		return nil, fmt.Errorf("新邮箱不能与当前邮箱相同")
	}

	// 验证邮箱验证码
	err = l.verifyEmailCode(req.Email, req.Code, "update_email")
	if err != nil {
		return nil, err
	}

	// 检查新邮箱是否已被其他用户使用
	var existingUser user_models.UserModel
	if err := l.svcCtx.DB.Where("email = ? AND user_id != ?", req.Email, req.UserID).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 更新用户邮箱
	err = l.svcCtx.DB.Model(&user).Update("email", req.Email).Error
	if err != nil {
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "用户修改邮箱成功",
		Data: map[string]interface{}{"userId": req.UserID},
	})
	return &types.UpdateEmailRes{}, nil
}

// 验证邮箱验证码
func (l *UpdateEmailLogic) verifyEmailCode(email, code, codeType string) error {
	// 从Redis获取存储的验证码
	codeKey := fmt.Sprintf("email_code_%s_%s", email, codeType)
	storedCode, err := l.svcCtx.Redis.Get(codeKey).Result()
	if err != nil {
		return fmt.Errorf("验证码已过期或不存在")
	}

	// 验证验证码
	if storedCode != code {
		return fmt.Errorf("验证码错误")
	}

	// 验证成功后删除验证码
	l.svcCtx.Redis.Del(codeKey)

	return nil
}
