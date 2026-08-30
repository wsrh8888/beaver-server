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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendVerifyDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendVerifyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVerifyDetailLogic {
	return &GetFriendVerifyDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetFriendVerifyDetail 管理后台：好友验证详情。
// admin 职责：校验 verifyId，复用 ListFriendVerifies 查单条，UserRpc 组装双方信息。
// RPC 职责：ListFriendVerifies(verify_id) 精确查询，不单独暴露 GetFriendVerify RPC。
func (l *GetFriendVerifyDetailLogic) GetFriendVerifyDetail(req *types.GetFriendVerifyDetailReq) (resp *types.GetFriendVerifyDetailRes, err error) {
	if req.VerifyID == "" {
		return nil, errors.New("验证记录ID不能为空")
	}

	rpcRes, err := l.svcCtx.FriendRpc.ListFriendVerifies(l.ctx, &friend_rpc.ListFriendVerifiesReq{
		VerifyId: req.VerifyID,
	})
	if err != nil {
		l.Errorf("获取好友验证详情失败: %v", err)
		return nil, err
	}
	if len(rpcRes.List) == 0 {
		return nil, errors.New("验证记录不存在")
	}

	v := rpcRes.List[0]
	users := map[string]*user_rpc.UserInfo{}
	if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: []string{v.SendUserId, v.RevUserId}}); err == nil && res != nil {
		users = res.UserInfo
	}

	sendName, sendAvatar, revName, revAvatar := "", "", "", ""
	if u, ok := users[v.SendUserId]; ok && u != nil {
		sendName, sendAvatar = u.NickName, u.Avatar
	}
	if u, ok := users[v.RevUserId]; ok && u != nil {
		revName, revAvatar = u.NickName, u.Avatar
	}

	return &types.GetFriendVerifyDetailRes{
		Id:               v.VerifyId,
		SendUserId:       v.SendUserId,
		SendUserName:     sendName,
		SendUserFileName: sendAvatar,
		RevUserId:        v.RevUserId,
		RevUserName:      revName,
		RevUserFileName:  revAvatar,
		SendStatus:       int(v.SendStatus),
		RevStatus:        int(v.RevStatus),
		Message:          v.Message,
		CreateTime:       v.CreatedAt,
		UpdateTime:       v.UpdatedAt,
	}, nil
}
