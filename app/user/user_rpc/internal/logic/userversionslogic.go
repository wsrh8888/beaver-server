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

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserVersionsLogic {
	return &UserVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserVersionsLogic) UserVersions(in *user_rpc.UserVersionsReq) (*user_rpc.UserVersionsRes, error) {
	// 对空数组进行处理
	if len(in.UserIds) == 0 {
		return &user_rpc.UserVersionsRes{
			UserVersions: make(map[string]int64),
		}, nil
	}

	// 查询指定用户ID列表的版本信息
	var userList []user_models.UserModel
	err := l.svcCtx.DB.Select("user_id, version").Where("user_id IN ?", in.UserIds).Find(&userList).Error

	if err != nil {
		l.Logger.Errorf("查询用户版本信息失败: %v", err)
		return nil, err
	}

	// 构造响应
	resp := &user_rpc.UserVersionsRes{
		UserVersions: make(map[string]int64, len(userList)),
	}

	for _, user := range userList {
		resp.UserVersions[user.UserID] = user.Version
	}

	// 为不存在的用户设置默认版本号0
	for _, userId := range in.UserIds {
		if _, exists := resp.UserVersions[userId]; !exists {
			resp.UserVersions[userId] = 0 // 用户不存在返回版本号0
		}
	}

	l.Logger.Infof("查询到 %d 个用户版本信息，请求ID数量: %d", len(userList), len(in.UserIds))
	return resp, nil
}
