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
	"encoding/json"
	"errors"

	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_api/internal/utils"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/open/openevent"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveRobotFromGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveRobotFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveRobotFromGroupLogic {
	return &RemoveRobotFromGroupLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RemoveRobotFromGroupLogic) RemoveRobotFromGroup(req *types.RemoveRobotFromGroupReq, authorization string) (resp *types.RemoveRobotFromGroupRes, err error) {
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
	if err := utils.RequireAppCapability(app, true, false); err != nil {
		return nil, err
	}

	robot, err := utils.EnsureAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, app)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.GroupRpc.RemoveGroupMember(l.ctx, &group_rpc.RemoveGroupMemberReq{
		GroupId:    req.GroupID,
		UserId:     robot.RobotID,
		OperatedBy: token.AppID,
	})
	if err != nil {
		return nil, err
	}

	go func() {
		body, _ := json.Marshal(map[string]interface{}{
			"group_id":    req.GroupID,
			"robot_id":    robot.RobotID,
			"operator_id": token.AppID,
		})
		_, _ = l.svcCtx.OpenRpc.DispatchPlatformEvent(context.Background(), &open_rpc.DispatchPlatformEventReq{
			AppId:     token.AppID,
			EventType: openevent.EventIMChatMemberBotRemoved,
			EventJson: string(body),
		})
	}()

	return &types.RemoveRobotFromGroupRes{
		Success: true,
	}, nil
}
