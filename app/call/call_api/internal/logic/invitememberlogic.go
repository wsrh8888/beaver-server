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

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type InviteMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 群聊中邀请成员入场
func NewInviteMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteMemberLogic {
	return &InviteMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("invite_member", ctx),
	}
}

func (l *InviteMemberLogic) InviteMember(req *types.InviteCallMemberReq) (resp *types.InviteCallMemberRes, err error) {
	// 1. 获取会话信息 (主要为了拿到 ConversationID 和 CallType)
	session, err := l.svcCtx.CallRpc.GetSession(l.ctx, &call_rpc.GetSessionReq{RoomId: req.RoomID})
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "获取通话会话失败",
			Data: map[string]any{"roomId": req.RoomID, "err": err.Error()},
		})
		return nil, err
	}

	// 2. 获取邀请者（发起人）信息
	callerInfo, _ := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: req.UserID})
	var callerUserInfo map[string]string
	if callerInfo != nil && callerInfo.GetUserInfo() != nil {
		callerUserInfo = map[string]string{
			"userId":   req.UserID,
			"nickName": callerInfo.GetUserInfo().NickName,
			"avatar":   callerInfo.GetUserInfo().Avatar,
		}
	}

	// 3. 依次登记状态并发送信令
	for _, targetID := range req.TargetIds {
		// 登记为待接听状态
		_, _ = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
			RoomId: req.RoomID,
			UserId: targetID,
			Status: 1, // 1-待接听 (ParticipantStatusCalling)
		})

		// 通过 RocketMQ 异步发送 WebSocket RTC_INVITE 信令给受邀方 (告知来电)
		payload := map[string]interface{}{
			"command":  wsCommandConst.CALL,
			"type":     wsTypeConst.CallReceive,
			"senderId": req.UserID,
			"targetId": targetID,
			"body": map[string]interface{}{
				"type":           call_models.SignalInvite,
				"roomId":         req.RoomID,
				"callerId":       req.UserID,
				"callType":       session.CallType,
				"callerUserInfo": callerUserInfo,
				"timestamp":      time.Now().Unix(),
			},
			"conversationId": session.ConversationId,
		}
		l.svcCtx.RocketMQ.SendMessage(l.ctx, mqwsconst.MqTopicWs, payload)

		// 4. [核心修复] 通知房间里的所有人：有新成员正在被呼叫中 (包括自己的其他设备同步)
		for _, pid := range session.ParticipantIds {
			if pid != targetID {
				payload := map[string]interface{}{
					"command":  wsCommandConst.CALL,
					"type":     wsTypeConst.CallReceive,
					"senderId": req.UserID,
					"targetId": pid,
					"body": map[string]interface{}{
						"type":   call_models.SignalInvite,
						"userId": targetID,
						"roomId": req.RoomID,
						"status": 1,
					},
					"conversationId": session.ConversationId,
				}
				l.svcCtx.RocketMQ.SendMessage(l.ctx, mqwsconst.MqTopicWs, payload)
			}
		}

		// 5. 开启超时处理定时器 (60秒未接听则自动设为超时)
		l.startTimeoutTimer(req.RoomID, targetID)
	}

	l.logger.Info(model.LogMsg{
		Text: "邀请通话成员成功",
		Data: map[string]interface{}{
			"roomId":         req.RoomID,
			"callerId":       req.UserID,
			"targetIds":      req.TargetIds,
			"conversationId": session.ConversationId,
		},
	})

	return &types.InviteCallMemberRes{}, nil
}

// startTimeoutTimer 异步计时器：如果用户在规定时间内未接听，则自动更新状态为超时
func (l *InviteMemberLogic) startTimeoutTimer(roomID, userID string) {
	time.AfterFunc(60*time.Second, func() {
		// 使用 context.Background()，因为原始请求的上下文会因接口返回而取消
		ctx := context.Background()

		// 1. 确认用户当前状态
		participants, err := l.svcCtx.CallRpc.GetParticipants(ctx, &call_rpc.GetParticipantsReq{RoomId: roomID})
		if err != nil {
			return
		}

		isStillCalling := false
		for _, p := range participants.Participants {
			if p.UserId == userID && p.Status == 1 { // 1 代表 Calling
				isStillCalling = true
				break
			}
		}

		// 2. 如果依然是待接听状态，则变更为超时
		if isStillCalling {
			_, _ = l.svcCtx.CallRpc.UpdateParticipantStatus(ctx, &call_rpc.UpdateParticipantStatusReq{
				RoomId: roomID,
				UserId: userID,
				Status: 4, // 4 代表 Timeout
			})
		}
	})
}
