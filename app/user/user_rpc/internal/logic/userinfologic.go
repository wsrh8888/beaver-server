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
	"time"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user_rpc.UserInfoReq) (*user_rpc.UserInfoRes, error) {
	var user user_models.UserModel

	err := l.svcCtx.DB.Take(&user, "user_id = ?", in.UserID).Error

	if err != nil {
		logx.Errorf("查询用户失败: %s", err.Error())
		return nil, errors.New("用户不存在")
	}

	return &user_rpc.UserInfoRes{UserInfo: &user_rpc.UserInfo{
		UserId:    user.UserID,
		NickName:  user.NickName,
		Avatar:    user.Avatar,
		Version:   user.Version,
		Email:     user.Email,
		Abstract:  user.Abstract,
		Phone:     user.Phone,
		Status:    int32(user.Status),
		Source:    user.Source,
		UserType:  int32(user.UserType),
		CreatedAt: time.Time(user.CreatedAt).Format(time.RFC3339),
		UpdatedAt: time.Time(user.UpdatedAt).Format(time.RFC3339),
	}}, nil
}
