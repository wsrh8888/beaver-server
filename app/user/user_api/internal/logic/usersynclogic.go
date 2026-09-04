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
	"strings"
	"time"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type UserSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 用户数据同步
func NewUserSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSyncLogic {
	return &UserSyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("user_sync", ctx),
	}
}

func (l *UserSyncLogic) UserSync(req *types.UserSyncReq) (resp *types.UserSyncRes, err error) {
	if len(req.UserVersions) == 0 {
		l.logger.Info(model.LogMsg{
			Text: "未指定同步用户数据",
			Data: map[string]any{"userId": req.UserID},
		})
		return &types.UserSyncRes{Users: []types.UserSyncItem{}}, nil
	}

	// 构建查询条件
	var conditions []string
	var args []interface{}

	for _, uv := range req.UserVersions {
		conditions = append(conditions, "(user_id = ? AND version >= ?)")
		args = append(args, uv.UserID, uv.Version)
	}

	// 查询并转换用户数据
	var users []user_models.UserModel
	if err = l.svcCtx.DB.Where(strings.Join(conditions, " OR "), args...).
		Order("version ASC").Find(&users).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "查询相关用户数据失败",
			Data: map[string]any{"userId": req.UserID, "err": err.Error()},
		})
		return nil, err
	}

	// 转换为响应格式
	userItems := make([]types.UserSyncItem, len(users))
	for i, user := range users {
		userItems[i] = types.UserSyncItem{
			UserID:    user.UserID,
			NickName:  user.NickName,
			Avatar:    user.Avatar,
			Abstract:  user.Abstract,
			Phone:     user.Phone,
			Email:     user.Email,
			Gender:    user.Gender,
			Status:    user.Status,
			UserType:  user.UserType,
			Version:   user.Version,
			CreatedAt: time.Time(user.CreatedAt).Unix(),
			UpdatedAt: time.Time(user.UpdatedAt).Unix(),
		}
	}

	l.logger.Info(model.LogMsg{
		Text: "用户数据同步完成",
		Data: map[string]interface{}{
			"userId":       req.UserID,
			"requestCount": len(req.UserVersions),
			"resultCount":  len(userItems),
		},
	})

	return &types.UserSyncRes{Users: userItems}, nil
}
