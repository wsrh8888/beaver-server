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
	"fmt"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/conversation"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchRes, err error) {
	// 参数验证
	if req.Keyword == "" {
		return nil, errors.New("搜索关键词不能为空")
	}

	var userInfo *user_rpc.UserInfo
	var userId string

	// 根据搜索类型查询用户信息
	switch req.Type {
	case "email":
		// 根据邮箱查询
		searchResp, err := l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
			Keyword: req.Keyword,
			Type:    "email",
		})
		if err != nil {
			l.Logger.Errorf("根据邮箱查询用户失败: email=%s, error=%v", req.Keyword, err)
			return nil, errors.New("用户不存在")
		}
		if searchResp.UserInfo == nil {
			l.Logger.Errorf("根据邮箱查询用户失败: email=%s, 未找到用户", req.Keyword)
			return nil, errors.New("用户不存在")
		}
		userInfo = searchResp.UserInfo
		userId = userInfo.UserId
	case "userId":
		// 根据用户ID查询
		searchResp, err := l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
			Keyword: req.Keyword,
			Type:    "userId",
		})
		if err != nil {
			l.Logger.Errorf("根据用户ID查询用户失败: userId=%s, error=%v", req.Keyword, err)
			return nil, errors.New("用户不存在")
		}
		if searchResp.UserInfo == nil {
			l.Logger.Errorf("根据用户ID查询用户失败: userId=%s, 未找到用户", req.Keyword)
			return nil, errors.New("用户不存在")
		}
		userInfo = searchResp.UserInfo
		userId = userInfo.UserId
	default:
		// 默认按邮箱搜索
		searchResp, err := l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
			Keyword: req.Keyword,
			Type:    "email",
		})
		if err != nil {
			l.Logger.Errorf("根据邮箱查询用户失败: email=%s, error=%v", req.Keyword, err)
			return nil, errors.New("用户不存在")
		}
		if searchResp.UserInfo == nil {
			l.Logger.Errorf("根据邮箱查询用户失败: email=%s, 未找到用户", req.Keyword)
			return nil, errors.New("用户不存在")
		}
		userInfo = searchResp.UserInfo
		userId = userInfo.UserId
	}

	// 获取完整的用户信息（包括Abstract字段）
	userInfoResp, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: userId,
	})
	if err != nil {
		l.Logger.Errorf("获取用户详细信息失败: userID=%s, error=%v", userId, err)
		return nil, errors.New("获取用户详细信息失败")
	}

	userDetail := userInfoResp.UserInfo

	// 不能搜索自己
	if req.UserID == userInfo.UserId {
		return nil, errors.New("不能搜索自己")
	}

	// 获取好友关系
	var friend friend_models.FriendModel
	isFriend := friend.IsFriend(l.svcCtx.DB, req.UserID, userInfo.UserId)

	// 生成会话Id
	conversationID, err := conversation.GenerateConversation([]string{req.UserID, userInfo.UserId})
	if err != nil {
		l.Logger.Errorf("生成会话Id失败: %v", err)
		return nil, fmt.Errorf("生成会话Id失败: %v", err)
	}

	// 构造返回值
	resp = &types.SearchRes{
		UserID:         userInfo.UserId,
		NickName:       userInfo.NickName,
		Avatar:         userInfo.Avatar,
		Abstract:       userDetail.Abstract,
		IsFriend:       isFriend,
		ConversationID: conversationID,
		Email:          userDetail.Email,
	}

	l.Logger.Infof("搜索用户成功: userID=%s, keyword=%s, type=%s, isFriend=%v", req.UserID, req.Keyword, req.Type, isFriend)
	return resp, nil
}
