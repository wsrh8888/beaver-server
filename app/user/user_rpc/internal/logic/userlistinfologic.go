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
	"time"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserListInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListInfoLogic {
	return &UserListInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserListInfoLogic) UserListInfo(in *user_rpc.UserListInfoReq) (*user_rpc.UserListInfoRes, error) {
	// 对空数组进行处理
	if len(in.UserIdList) == 0 {
		return &user_rpc.UserListInfoRes{
			UserInfo: make(map[string]*user_rpc.UserInfo),
		}, nil
	}

	// 构建查询条件
	query := l.svcCtx.DB.Model(&user_models.UserModel{}).Where("user_id IN ?", in.UserIdList)

	// 如果提供了时间戳，则只返回该时间之后更新的用户
	if in.SinceTimestamp > 0 {
		query = query.Where("updated_at > ?", in.SinceTimestamp)
	}

	var userList []user_models.UserModel
	err := query.Find(&userList).Error

	if err != nil {
		l.Logger.Errorf("查询用户列表失败: %v", err)
		return nil, err
	}

	resp := &user_rpc.UserListInfoRes{
		UserInfo: make(map[string]*user_rpc.UserInfo, len(userList)),
	}

	for _, user := range userList {
		resp.UserInfo[user.UserID] = &user_rpc.UserInfo{
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
		}
	}

	if in.SinceTimestamp > 0 {
		l.Logger.Infof("增量查询用户，时间戳: %d，查询到 %d 个用户信息，请求ID数量: %d",
			in.SinceTimestamp, len(userList), len(in.UserIdList))
	} else {
		l.Logger.Infof("全量查询用户，查询到 %d 个用户信息，请求ID数量: %d",
			len(userList), len(in.UserIdList))
	}

	return resp, nil
}
