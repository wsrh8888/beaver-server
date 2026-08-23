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

package call_models

/*
RTC 实时信令事件说明 (用于 WebSocket 瞬时通知，不进入持久化聊天记录表)：
1. RTC_INVITE  (呼叫): 发起通话，被叫方收到此信号显示振铃界面。
2. RTC_CANCEL  (取消): 发起方在对方接听前主动挂断，被叫方停止振铃。
3. RTC_ACCEPT  (接听): 被叫方接听，发起方收到后进入通话状态。
4. RTC_REJECT  (拒绝): 被叫方拒绝接听，发起方显示对方拒绝并退出。
5. RTC_HANGUP  (挂断): 通话过程中任意一方挂断，所有人退出房间。
*/

const (
	SignalInvite = "RTC_INVITE"   // 呼叫邀请
	SignalCancel = "RTC_CANCEL"   // 取消呼叫
	SignalAccept = "RTC_ACCEPTED" // 接听通话
	SignalReject = "RTC_REJECT"   // 拒绝接听
	SignalHangup = "RTC_HANGUP"   // 通话中挂断
)
