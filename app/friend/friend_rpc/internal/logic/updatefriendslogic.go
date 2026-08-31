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

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

const (
	friendActionHardDelete int32 = 1 // 物理删除好友关系
	friendActionRestore    int32 = 2 // 恢复软删除的好友关系
)

type UpdateFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendsLogic {
	return &UpdateFriendsLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("update_friends", ctx)}
}

func (l *UpdateFriendsLogic) UpdateFriends(in *friend_rpc.UpdateFriendsReq) (*friend_rpc.UpdateFriendsRes, error) {
	switch in.Action {
	case friendActionHardDelete:
		return l.hardDelete(in.RelationIds)
	case friendActionRestore:
		return l.restore(in.RelationIds)
	default:
		return nil, errors.New("不支持的操作类型")
	}
}

func (l *UpdateFriendsLogic) hardDelete(relationIDs []string) (*friend_rpc.UpdateFriendsRes, error) {
	var ids []uint
	for _, rid := range relationIDs {
		f, err := findFriend(l.svcCtx.DB, rid)
		if err != nil {
			continue
		}
		ids = append(ids, f.Id)
	}
	if len(ids) == 0 {
		return &friend_rpc.UpdateFriendsRes{}, nil
	}
	if err := l.svcCtx.DB.Unscoped().Where("id IN ?", ids).Delete(&friend_models.FriendModel{}).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "删除好友失败", Data: map[string]any{"ids": ids, "err": err.Error()}})
		return nil, err
	}
	l.logger.Info(model.LogMsg{Text: "删除好友成功", Data: map[string]interface{}{"count": len(ids)}})
	return &friend_rpc.UpdateFriendsRes{AffectedCount: int64(len(ids))}, nil
}

func (l *UpdateFriendsLogic) restore(relationIDs []string) (*friend_rpc.UpdateFriendsRes, error) {
	var affected int64
	for _, rid := range relationIDs {
		f, err := findFriend(l.svcCtx.DB, rid)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "查询好友关系失败", Data: map[string]any{"relationId": rid, "err": err.Error()}})
			return nil, err
		}
		if !f.IsDeleted {
			affected++
			continue
		}
		if err := l.svcCtx.DB.Model(f).Update("is_deleted", false).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "恢复好友失败", Data: map[string]any{"relationId": rid, "err": err.Error()}})
			return nil, err
		}
		affected++
	}
	l.logger.Info(model.LogMsg{Text: "恢复好友成功", Data: map[string]interface{}{"affected": affected}})
	return &friend_rpc.UpdateFriendsRes{AffectedCount: affected}, nil
}
