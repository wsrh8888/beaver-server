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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type DeleteUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUsersLogic {
	return &DeleteUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("delete_users", ctx),
	}
}

func (l *DeleteUsersLogic) DeleteUsers(in *user_rpc.DeleteUsersReq) (*user_rpc.DeleteUsersRes, error) {
	if len(in.UserIds) == 0 {
		return &user_rpc.DeleteUsersRes{}, nil
	}
	result := l.svcCtx.DB.Model(&user_models.UserModel{}).
		Where("user_id IN ?", in.UserIds).
		Update("status", 3)
	if result.Error != nil {
		l.logger.Error(model.LogMsg{
			Text: "删除用户失败",
			Data: map[string]any{"userIds": in.UserIds, "err": result.Error.Error()},
		})
		return nil, result.Error
	}
	l.logger.Info(model.LogMsg{
		Text: "删除用户成功",
		Data: map[string]interface{}{"userIds": in.UserIds, "affected": result.RowsAffected},
	})
	return &user_rpc.DeleteUsersRes{AffectedCount: result.RowsAffected}, nil
}
