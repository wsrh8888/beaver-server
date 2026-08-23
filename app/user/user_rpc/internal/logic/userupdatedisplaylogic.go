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

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserUpdateDisplayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserUpdateDisplayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateDisplayLogic {
	return &UserUpdateDisplayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserUpdateDisplayLogic) UserUpdateDisplay(in *user_rpc.UserUpdateDisplayReq) (*user_rpc.UserUpdateDisplayRes, error) {
	if in.UserId == "" {
		return nil, errors.New("userId 不能为空")
	}
	if in.NickName == "" && in.Avatar == "" {
		return &user_rpc.UserUpdateDisplayRes{}, nil
	}

	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ?", in.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if in.NickName != "" {
		updates["nick_name"] = in.NickName
	}
	if in.Avatar != "" {
		updates["avatar"] = in.Avatar
	}

	version := l.svcCtx.VersionGen.GetNextVersion("users", "user_id", in.UserId)
	if version == -1 {
		return nil, errors.New("获取用户版本号失败")
	}
	updates["version"] = version

	if err := l.svcCtx.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user_rpc.UserUpdateDisplayRes{Version: version}, nil
}
