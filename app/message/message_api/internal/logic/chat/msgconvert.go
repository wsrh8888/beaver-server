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
	"beaver/app/message/message_api/internal/types"
	chat_rpc "beaver/app/message/message_rpc/pb"
)

// pbMsgToPeerMsg 把 rpc 的消息结构显式映射到 API 出参。
// 消息体含自引用（WithdrawMsg/ReplyMsg 里嵌 Msg），故递归转换，不能用反射盲配
// （那会在自引用上指数展开）。字段逐一写出：pb 的 protoc 命名（Id/Url）→ API 的 Go 规范命名（ID/URL）。
func pbMsgToPeerMsg(m *chat_rpc.Msg) *types.PeerMsg {
	if m == nil {
		return nil
	}
	out := &types.PeerMsg{
		Type:        m.Type,
		TargetMsgID: m.TargetMsgId,
		AtUserIDs:   m.AtUserIds,
	}
	if v := m.GetTextMsg(); v != nil {
		out.TextMsg = &types.PeerTextMsg{Content: v.Content}
	}
	if v := m.GetImageMsg(); v != nil {
		out.ImageMsg = &types.PeerImageMsg{FileUrl: v.FileUrl, Width: v.Width, Height: v.Height, Size: v.Size}
	}
	if v := m.GetVideoMsg(); v != nil {
		out.VideoMsg = &types.PeerVideoMsg{FileUrl: v.FileUrl, Width: v.Width, Height: v.Height, Duration: v.Duration, ThumbnailUrl: v.ThumbnailUrl, Size: v.Size}
	}
	if v := m.GetFileMsg(); v != nil {
		out.FileMsg = &types.PeerFileMsg{FileUrl: v.FileUrl, FileName: v.FileName, Size: v.Size, MimeType: v.MimeType, Extension: v.Extension, OpenMode: v.OpenMode}
	}
	if v := m.GetVoiceMsg(); v != nil {
		out.VoiceMsg = &types.PeerVoiceMsg{FileUrl: v.FileUrl, Duration: v.Duration, Size: v.Size}
	}
	if v := m.GetEmojiMsg(); v != nil {
		out.EmojiMsg = &types.PeerEmojiMsg{FileUrl: v.FileUrl, EmojiID: v.EmojiId, PackageID: v.PackageId, Width: v.Width, Height: v.Height}
	}
	if v := m.GetNotificationMsg(); v != nil {
		out.NotificationMsg = &types.PeerNotificationMsg{Type: v.Type, Actors: v.Actors}
	}
	if v := m.GetAudioFileMsg(); v != nil {
		out.AudioFileMsg = &types.PeerAudioFileMsg{FileUrl: v.FileUrl, FileName: v.FileName, Duration: v.Duration, Size: v.Size}
	}
	if v := m.GetCallMsg(); v != nil {
		out.CallMsg = &types.PeerCallMsg{RoomID: v.RoomId, CallType: v.CallType, Status: v.Status, Duration: v.Duration}
	}
	if v := m.GetWithdrawMsg(); v != nil {
		out.WithdrawMsg = &types.PeerWithdrawMsg{OriginMsgID: v.OriginMsgId, OriginMsg: pbMsgToPeerMsg(v.OriginMsg)}
	}
	if v := m.GetReplyMsg(); v != nil {
		out.ReplyMsg = &types.PeerReplyMsg{
			OriginMsgID: v.OriginMsgId,
			OriginMsg:   pbMsgToPeerMsg(v.OriginMsg),
			ReplyMsg:    pbMsgToPeerMsg(v.ReplyMsg),
		}
	}
	if v := m.GetForwardMsg(); v != nil {
		out.ForwardMsg = &types.PeerForwardMsg{Title: v.Title, RecordID: v.RecordId, Count: v.Count}
	}
	if v := m.GetMarkdownMsg(); v != nil {
		out.MarkdownMsg = &types.PeerMarkdownMsg{Content: v.Content, Title: v.Title}
	}
	if v := m.GetLinkMsg(); v != nil {
		out.LinkMsg = &types.PeerLinkMsg{Url: v.Url, Title: v.Title, Desc: v.Desc, ImageUrl: v.ImageUrl}
	}
	if v := m.GetCloudDocMsg(); v != nil {
		out.CloudDocMsg = &types.PeerCloudDocMsg{
			DocID: v.DocId, DocType: v.DocType, Title: v.Title, OwnerID: v.OwnerId,
			Perm: v.Perm, CoverUrl: v.CoverUrl, Revision: v.Revision,
		}
	}
	if v := m.GetCardMsg(); v != nil {
		out.CardMsg = &types.PeerCardMsg{CardType: v.CardType, ID: v.Id, ExpireAt: v.ExpireAt, InviteToken: v.InviteToken}
	}
	return out
}

// peerMsgToPbMsg 入站方向：API 请求里的 Msg 显式映射到 rpc 入参。
func peerMsgToPbMsg(m *types.PeerMsg) *chat_rpc.Msg {
	if m == nil {
		return nil
	}
	out := &chat_rpc.Msg{
		Type:        m.Type,
		TargetMsgId: m.TargetMsgID,
		AtUserIds:   m.AtUserIDs,
	}
	if v := m.TextMsg; v != nil {
		out.TextMsg = &chat_rpc.TextMsg{Content: v.Content}
	}
	if v := m.ImageMsg; v != nil {
		out.ImageMsg = &chat_rpc.ImageMsg{FileUrl: v.FileUrl, Width: v.Width, Height: v.Height, Size: v.Size}
	}
	if v := m.VideoMsg; v != nil {
		out.VideoMsg = &chat_rpc.VideoMsg{FileUrl: v.FileUrl, Width: v.Width, Height: v.Height, Duration: v.Duration, ThumbnailUrl: v.ThumbnailUrl, Size: v.Size}
	}
	if v := m.FileMsg; v != nil {
		out.FileMsg = &chat_rpc.FileMsg{FileUrl: v.FileUrl, FileName: v.FileName, Size: v.Size, MimeType: v.MimeType, Extension: v.Extension, OpenMode: v.OpenMode}
	}
	if v := m.VoiceMsg; v != nil {
		out.VoiceMsg = &chat_rpc.VoiceMsg{FileUrl: v.FileUrl, Duration: v.Duration, Size: v.Size}
	}
	if v := m.EmojiMsg; v != nil {
		out.EmojiMsg = &chat_rpc.EmojiMsg{FileUrl: v.FileUrl, EmojiId: v.EmojiID, PackageId: v.PackageID, Width: v.Width, Height: v.Height}
	}
	if v := m.NotificationMsg; v != nil {
		out.NotificationMsg = &chat_rpc.NotificationMsg{Type: v.Type, Actors: v.Actors}
	}
	if v := m.AudioFileMsg; v != nil {
		out.AudioFileMsg = &chat_rpc.AudioFileMsg{FileUrl: v.FileUrl, FileName: v.FileName, Duration: v.Duration, Size: v.Size}
	}
	if v := m.CallMsg; v != nil {
		out.CallMsg = &chat_rpc.CallMsg{RoomId: v.RoomID, CallType: v.CallType, Status: v.Status, Duration: v.Duration}
	}
	if v := m.WithdrawMsg; v != nil {
		out.WithdrawMsg = &chat_rpc.WithdrawMsg{OriginMsgId: v.OriginMsgID, OriginMsg: peerMsgToPbMsg(v.OriginMsg)}
	}
	if v := m.ReplyMsg; v != nil {
		out.ReplyMsg = &chat_rpc.ReplyMsg{
			OriginMsgId: v.OriginMsgID,
			OriginMsg:   peerMsgToPbMsg(v.OriginMsg),
			ReplyMsg:    peerMsgToPbMsg(v.ReplyMsg),
		}
	}
	if v := m.ForwardMsg; v != nil {
		out.ForwardMsg = &chat_rpc.ForwardMsg{Title: v.Title, RecordId: v.RecordID, Count: v.Count}
	}
	if v := m.MarkdownMsg; v != nil {
		out.MarkdownMsg = &chat_rpc.MarkdownMsg{Content: v.Content, Title: v.Title}
	}
	if v := m.LinkMsg; v != nil {
		out.LinkMsg = &chat_rpc.LinkMsg{Url: v.Url, Title: v.Title, Desc: v.Desc, ImageUrl: v.ImageUrl}
	}
	if v := m.CloudDocMsg; v != nil {
		out.CloudDocMsg = &chat_rpc.CloudDocMsg{
			DocId: v.DocID, DocType: v.DocType, Title: v.Title, OwnerId: v.OwnerID,
			Perm: v.Perm, CoverUrl: v.CoverUrl, Revision: v.Revision,
		}
	}
	if v := m.CardMsg; v != nil {
		out.CardMsg = &chat_rpc.CardMsg{CardType: v.CardType, Id: v.ID, ExpireAt: v.ExpireAt, InviteToken: v.InviteToken}
	}
	return out
}

func pbSenderToPeerSender(s *chat_rpc.Sender) *types.PeerSender {
	if s == nil {
		return nil
	}
	return &types.PeerSender{UserID: s.UserId, Avatar: s.Avatar, NickName: s.NickName, UserType: s.UserType}
}
