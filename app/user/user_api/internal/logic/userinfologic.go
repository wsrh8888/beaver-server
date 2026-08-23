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

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoRes, err error) {
	fmt.Printf("获取用户的基础信息, UserID: %v\n", req.UserID)

	// 直接从数据库查询，避免RPC调用自身服务
	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "user_id = ?", req.UserID).Error
	if err != nil {
		fmt.Printf("[ERROR] 查询用户失败, UserID: %v, error: %v\n", req.UserID, err)
		return nil, err
	}

	resp = &types.UserInfoRes{
		UserID:   user.UserID,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Abstract: user.Abstract,
		Phone:    user.Phone,
		Email:    user.Email,
		Gender:   user.Gender,
		Version:  user.Version,
	}

	return resp, nil
}
