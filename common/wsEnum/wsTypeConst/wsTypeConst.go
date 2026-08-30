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

package wsTypeConst

type Type string

// send（发起消息给服务端）
// receive（发给对方设备）
// sync（发给自己其他设备进行记录同步）

const (
	PrivateMessageSend Type = "private_message_send" // 客户端->服务端 私聊消息发送
	GroupMessageSend   Type = "group_message_send"   // 客户端->服务端 群聊消息发送
	// --------------------------------------------------------
	// --------------------------------------------------------

	// 会话信息同步
	ChatConversationMetaReceive    Type = "chat_conversation_meta_receive"    //  服务端->客户端 会话信息同步
	ChatUserConversationReceive    Type = "chat_user_conversation_receive"    //  服务端->客户端 用户会话信息同步
	ChatConversationMessageReceive Type = "chat_conversation_message_receive" //  服务端->客户端 会话消息同步
	ChatPeerReadReceive            Type = "chat_peer_read_receive"            //  服务端->客户端 对端已读序列号同步
	ChatMessageMediaReceive        Type = "chat_message_media_receive"        //  服务端->客户端 消息媒体状态同步（语音已听等）
	TypingSend                     Type = "typing_send"                       //  客户端->服务端 正在输入
	TypingReceive                  Type = "typing_receive"                      //  服务端->客户端 对端正在输入
)
const (
	// -------------------------------------------------------------------------------------
	FriendReceive       Type = "friend_receive"        // 服务端->客户端 好友信息同步
	FriendVerifyReceive Type = "friend_verify_receive" // 服务端->客户端 好友验证信息同步
)

// -------------------------------------------------------------------------------------

const (
	GroupReceive            Type = "group_receive"              // 服务端->客户端 群组信息同步
	GroupJoinRequestReceive Type = "group_join_request_receive" // 服务端->客户端 群成员添加请求
	GroupMemberReceive      Type = "group_member_receive"       // 服务端->客户端 群成员变动（加入，离开、被踢出等）

)

const (
	CircleReceive Type = "circle_receive" // 服务端->客户端 圈子信息同步
)

const (
	// --------------------------------------------------------
	UserReceive     Type = "user_receive"      // 服务端->客户端 用户信息同步
	UserKickReceive Type = "user_kick_receive" // 服务端->客户端 设备被强制下线（携带deviceId，客户端比对后执行本地登出）
)

const (
	// 通知中心
	NotificationReceive         Type = "notification_receive"           // 服务端->客户端 通知推送
	NotificationMarkReadReceive Type = "notification_mark_read_receive" // 服务端->客户端 标记已读同步
)

const (
	// 表情中心
	EmojiReceive Type = "emoji_receive" // 服务端->客户端 表情数据同步
)

const (
	// 音视频通话
	CallReceive Type = "call_receive" // 服务端->客户端 通话信令同步（LiveKit 来电通知）
)
