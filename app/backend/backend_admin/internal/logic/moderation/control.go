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

package moderation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"encoding/json"
)

const (
	chatMessageStatusDeleted int32 = 4
)

func executeControlAction(ctx context.Context, svcCtx *svc.ServiceContext, operatorID string, caseID uint64, act types.ModerationControlAction) error {
	target := strings.TrimSpace(act.Target)
	extra := strings.TrimSpace(act.Extra)
	detail, _ := json.Marshal(act)

	var err error
	switch act.Action {
	case "ban_user":
		if target == "" {
			return errors.New("封禁用户需要指定 userId")
		}
		_, err = svcCtx.UserRpc.UpdateUsersStatus(ctx, &user_rpc.UpdateUsersStatusReq{
			UserIds: []string{target},
			Status:  2,
		})
	case "unban_user":
		if target == "" {
			return errors.New("解封用户需要指定 userId")
		}
		_, err = svcCtx.UserRpc.UpdateUsersStatus(ctx, &user_rpc.UpdateUsersStatusReq{
			UserIds: []string{target},
			Status:  1,
		})
	case "delete_message":
		if target == "" {
			return errors.New("删除消息需要指定 messageId")
		}
		_, err = svcCtx.ChatRpc.UpdateChatMessages(ctx, &chat_rpc.UpdateChatMessagesReq{
			MessageIds: []string{target},
			Status:     chatMessageStatusDeleted,
		})
	case "clear_conversation":
		if target == "" {
			return errors.New("清空会话需要指定 conversationId")
		}
		_, err = svcCtx.ChatRpc.UpdateChatMessages(ctx, &chat_rpc.UpdateChatMessagesReq{
			ConversationId: target,
			Status:         chatMessageStatusDeleted,
		})
	case "delete_moment":
		if target == "" {
			return errors.New("删除动态需要指定 momentId")
		}
		deleted := true
		_, err = svcCtx.MomentRpc.UpdateMoment(ctx, &moment_rpc.UpdateMomentReq{
			MomentId:  target,
			IsDeleted: &deleted,
		})
	case "dissolve_group":
		if target == "" {
			return errors.New("解散群组需要指定 groupId")
		}
		res, listErr := svcCtx.GroupRpc.ListGroups(ctx, &group_rpc.ListGroupsReq{
			GroupId:  target,
			Page:     1,
			PageSize: 1,
		})
		if listErr != nil {
			return listErr
		}
		if len(res.List) == 0 {
			return errors.New("群组不存在")
		}
		_, err = svcCtx.GroupRpc.UpdateGroup(ctx, &group_rpc.UpdateGroupReq{
			Id:     res.List[0].Id,
			Status: 3,
		})
	case "kick_member":
		groupID := extra
		userID := target
		if groupID == "" || userID == "" {
			return errors.New("踢出成员需要 extra=groupId 且 target=userId")
		}
		_, err = svcCtx.GroupRpc.RemoveGroupMember(ctx, &group_rpc.RemoveGroupMemberReq{
			GroupId: groupID,
			UserId:  userID,
			Kick:    true,
		})
	default:
		return fmt.Errorf("不支持的管控动作: %s", act.Action)
	}

	result := "success"
	errMsg := ""
	if err != nil {
		result = "failed"
		errMsg = err.Error()
	}
	svcCtx.RecordOperation(operatorID, act.Action, "control", target, caseID, string(detail), result, errMsg)
	return err
}
