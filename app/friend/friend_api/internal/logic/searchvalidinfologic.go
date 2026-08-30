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

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SearchValidInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchValidInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchValidInfoLogic {
	return &SearchValidInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchValidInfoLogic) SearchValidInfo(req *types.SearchValidInfoReq) (resp *types.SearchValidInfoRes, err error) {
	// 参数验证
	if req.UserID == "" || req.FriendID == "" {
		return nil, errors.New("用户ID和好友ID不能为空")
	}

	// 不能查询自己
	if req.UserID == req.FriendID {
		return nil, errors.New("不能查询自己")
	}

	var friendVerify friend_models.FriendVerifyModel

	// 查询好友验证记录
	err = l.svcCtx.DB.Where(
		"(rev_user_id = ? and send_user_id = ?) or (rev_user_id = ? and send_user_id = ?)",
		req.UserID, req.FriendID, req.FriendID, req.UserID,
	).First(&friendVerify).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Infof("好友验证记录不存在: userID=%s, friendID=%s", req.UserID, req.FriendID)
			return nil, errors.New("好友验证不存在")
		}
		l.Logger.Errorf("查询好友验证记录失败: %v", err)
		return nil, errors.New("查询好友验证记录失败")
	}

	// 检查验证状态
	if friendVerify.RevStatus != 0 {
		l.Logger.Errorf("好友验证已处理: verifyID=%s, status=%d", friendVerify.VerifyID, friendVerify.RevStatus)
		return nil, errors.New("该验证已处理，无法重复操作")
	}

	// 填充返回结果
	resp = &types.SearchValidInfoRes{
		ValidID: friendVerify.VerifyID, // 返回验证记录ID
	}

	l.Logger.Infof("查询好友验证信息成功: userID=%s, friendID=%s, validID=%s", req.UserID, req.FriendID, friendVerify.VerifyID)
	return resp, nil
}
