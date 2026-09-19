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

package robot

import (
	"context"
	"errors"

	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_api/internal/utils"
	"beaver/app/open/openevent"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type AddRobotToGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewAddRobotToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRobotToGroupLogic {
	return &AddRobotToGroupLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("add_robot_to_group", ctx)}
}

func (l *AddRobotToGroupLogic) AddRobotToGroup(req *types.AddRobotToGroupReq, authorization string) (resp *types.AddRobotToGroupRes, err error) {
	if req.GroupID == "" {
		return nil, errors.New("groupId 不能为空")
	}

	token, err := utils.ValidateAppAccessToken(l.svcCtx.DB, authorization)
	if err != nil {
		return nil, err
	}
	app, err := utils.LoadAppByID(l.svcCtx.DB, token.AppID)
	if err != nil {
		return nil, err
	}
	if err := utils.RequireAppEnabled(app); err != nil {
		return nil, err
	}

	robot, err := utils.EnsureAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, app)
	if err != nil {
		return nil, err
	}

	groupRes, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{
		GroupIDs: []string{req.GroupID},
	})
	if err != nil || len(groupRes.Groups) == 0 {
		return nil, errors.New("群组不存在")
	}

	_, err = l.svcCtx.GroupRpc.AddGroupMember(l.ctx, &group_rpc.AddGroupMemberReq{
		GroupId:    req.GroupID,
		UserId:     robot.RobotID,
		OperatedBy: token.AppID,
	})
	if err != nil {
		return nil, err
	}

	go func() {
		body := map[string]interface{}{
			"group_id":    req.GroupID,
			"robot_id":    robot.RobotID,
			"operator_id": token.AppID,
		}
		if err := openevent.Push(context.Background(), l.svcCtx.RocketMQ, robot.RobotID,
			openevent.EventIMChatMemberBotAdded, "group_"+req.GroupID, body); err != nil {
			l.logger.Error(model.LogMsg{
				Text: "机器人进群事件投递失败",
				Data: map[string]any{"robotId": robot.RobotID, "groupId": req.GroupID, "err": err.Error()},
			})
		}
	}()

	return &types.AddRobotToGroupRes{
		RobotID: robot.RobotID,
	}, nil
}
