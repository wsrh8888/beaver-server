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

package chat_message

import (
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"encoding/json"
)

func fileUrlFromMap(m map[string]interface{}) string {
	if fileUrl, ok := m["fileUrl"].(string); ok {
		return fileUrl
	}
	return ""
}

func thumbnailUrlFromMap(m map[string]interface{}) string {
	if thumbnailUrl, ok := m["thumbnailUrl"].(string); ok {
		return thumbnailUrl
	}
	return ""
}

// convertToRpcMsg 将原始消息转换为RPC消息格式
func convertToRpcMsg(msg json.RawMessage) (*chat_rpc.Msg, error) {
	var msgData map[string]interface{}
	err := json.Unmarshal(msg, &msgData)
	if err != nil {
		return nil, err
	}

	rpcMsg := &chat_rpc.Msg{}

	if msgType, ok := msgData["type"].(float64); ok {
		rpcMsg.Type = uint32(msgType)
	}

	switch rpcMsg.Type {
	case 1:
		if textMsg, ok := msgData["textMsg"].(map[string]interface{}); ok {
			if content, ok := textMsg["content"].(string); ok {
				rpcMsg.TextMsg = &chat_rpc.TextMsg{Content: content}
			}
		}
	case 2:
		if imageMsg, ok := msgData["imageMsg"].(map[string]interface{}); ok {
			rpcMsg.ImageMsg = &chat_rpc.ImageMsg{
				FileUrl: fileUrlFromMap(imageMsg),
			}
			if width, ok := imageMsg["width"].(float64); ok {
				rpcMsg.ImageMsg.Width = int32(width)
			}
			if height, ok := imageMsg["height"].(float64); ok {
				rpcMsg.ImageMsg.Height = int32(height)
			}
			if size, ok := imageMsg["size"].(float64); ok {
				rpcMsg.ImageMsg.Size = int64(size)
			}
		}
	case 3:
		if videoMsg, ok := msgData["videoMsg"].(map[string]interface{}); ok {
			rpcMsg.VideoMsg = &chat_rpc.VideoMsg{
				FileUrl:       fileUrlFromMap(videoMsg),
				ThumbnailUrl: thumbnailUrlFromMap(videoMsg),
			}
			if width, ok := videoMsg["width"].(float64); ok {
				rpcMsg.VideoMsg.Width = int32(width)
			}
			if height, ok := videoMsg["height"].(float64); ok {
				rpcMsg.VideoMsg.Height = int32(height)
			}
			if duration, ok := videoMsg["duration"].(float64); ok {
				rpcMsg.VideoMsg.Duration = int32(duration)
			}
			if size, ok := videoMsg["size"].(float64); ok {
				rpcMsg.VideoMsg.Size = int64(size)
			}
		}
	case 4:
		if fileMsg, ok := msgData["fileMsg"].(map[string]interface{}); ok {
			file := &chat_rpc.FileMsg{
				FileUrl: fileUrlFromMap(fileMsg),
			}
			if fileName, ok := fileMsg["fileName"].(string); ok {
				file.FileName = fileName
			}
			if size, ok := fileMsg["size"].(float64); ok {
				file.Size = int64(size)
			}
			if mimeType, ok := fileMsg["mimeType"].(string); ok {
				file.MimeType = mimeType
			}
			rpcMsg.FileMsg = file
		}
	case 5:
		if voiceMsg, ok := msgData["voiceMsg"].(map[string]interface{}); ok {
			rpcMsg.VoiceMsg = &chat_rpc.VoiceMsg{
				FileUrl: fileUrlFromMap(voiceMsg),
			}
			if duration, ok := voiceMsg["duration"].(float64); ok {
				rpcMsg.VoiceMsg.Duration = int32(duration)
			}
			if size, ok := voiceMsg["size"].(float64); ok {
				rpcMsg.VoiceMsg.Size = int64(size)
			}
		}
	case 6:
		if emojiMsg, ok := msgData["emojiMsg"].(map[string]interface{}); ok {
			emoji := &chat_rpc.EmojiMsg{
				FileUrl: fileUrlFromMap(emojiMsg),
			}
			if emojiId, ok := emojiMsg["emojiId"].(string); ok {
				emoji.EmojiId = emojiId
			}
			if packageId, ok := emojiMsg["packageId"].(string); ok {
				emoji.PackageId = packageId
			}
			if width, ok := emojiMsg["width"].(float64); ok {
				emoji.Width = int64(width)
			}
			if height, ok := emojiMsg["height"].(float64); ok {
				emoji.Height = int64(height)
			}
			rpcMsg.EmojiMsg = emoji
		}
	case 7:
		if notificationMsg, ok := msgData["notificationMsg"].(map[string]interface{}); ok {
			rpcMsg.NotificationMsg = &chat_rpc.NotificationMsg{}
			if msgType, ok := notificationMsg["type"].(float64); ok {
				rpcMsg.NotificationMsg.Type = int32(msgType)
			}
			if actors, ok := notificationMsg["actors"].([]interface{}); ok {
				for _, actor := range actors {
					if actorStr, ok := actor.(string); ok {
						rpcMsg.NotificationMsg.Actors = append(rpcMsg.NotificationMsg.Actors, actorStr)
					}
				}
			}
		}
	case 8:
		if audioFileMsg, ok := msgData["audioFileMsg"].(map[string]interface{}); ok {
			rpcMsg.AudioFileMsg = &chat_rpc.AudioFileMsg{
				FileUrl: fileUrlFromMap(audioFileMsg),
			}
			if fileName, ok := audioFileMsg["fileName"].(string); ok {
				rpcMsg.AudioFileMsg.FileName = fileName
			}
			if duration, ok := audioFileMsg["duration"].(float64); ok {
				rpcMsg.AudioFileMsg.Duration = int32(duration)
			}
			if size, ok := audioFileMsg["size"].(float64); ok {
				rpcMsg.AudioFileMsg.Size = int64(size)
			}
		}
	case 9:
		if callMsg, ok := msgData["callMsg"].(map[string]interface{}); ok {
			rpcMsg.CallMsg = &chat_rpc.CallMsg{}
			if roomId, ok := callMsg["roomId"].(string); ok {
				rpcMsg.CallMsg.RoomId = roomId
			}
			if callType, ok := callMsg["callType"].(float64); ok {
				rpcMsg.CallMsg.CallType = int32(callType)
			}
			if status, ok := callMsg["status"].(float64); ok {
				rpcMsg.CallMsg.Status = int32(status)
			}
			if duration, ok := callMsg["duration"].(float64); ok {
				rpcMsg.CallMsg.Duration = int64(duration)
			}
		}
	case 10:
		if withdrawMsg, ok := msgData["withdrawMsg"].(map[string]interface{}); ok {
			rpcMsg.WithdrawMsg = &chat_rpc.WithdrawMsg{}
			if originMsgId, ok := withdrawMsg["originMsgId"].(string); ok {
				rpcMsg.WithdrawMsg.OriginMsgId = originMsgId
			}
			if originMsgMap, ok := withdrawMsg["originMsg"].(map[string]interface{}); ok {
				jsonData, _ := json.Marshal(originMsgMap)
				originMsg, _ := convertToRpcMsg(jsonData)
				rpcMsg.WithdrawMsg.OriginMsg = originMsg
			}
		}
	case 11:
		if replyMsg, ok := msgData["replyMsg"].(map[string]interface{}); ok {
			rpcMsg.ReplyMsg = &chat_rpc.ReplyMsg{}
			if originMsgId, ok := replyMsg["originMsgId"].(string); ok {
				rpcMsg.ReplyMsg.OriginMsgId = originMsgId
			}
			if originMsgMap, ok := replyMsg["originMsg"].(map[string]interface{}); ok {
				jsonData, _ := json.Marshal(originMsgMap)
				originMsg, _ := convertToRpcMsg(jsonData)
				rpcMsg.ReplyMsg.OriginMsg = originMsg
			}
			if replyMsgInnerMap, ok := replyMsg["replyMsg"].(map[string]interface{}); ok {
				jsonData, _ := json.Marshal(replyMsgInnerMap)
				replyInnerMsg, _ := convertToRpcMsg(jsonData)
				rpcMsg.ReplyMsg.ReplyMsg = replyInnerMsg
			}
		}
	case 12:
		if forwardMsg, ok := msgData["forwardMsg"].(map[string]interface{}); ok {
			rpcMsg.ForwardMsg = &chat_rpc.ForwardMsg{}
			if title, ok := forwardMsg["title"].(string); ok {
				rpcMsg.ForwardMsg.Title = title
			}
			if recordId, ok := forwardMsg["recordId"].(string); ok {
				rpcMsg.ForwardMsg.RecordId = recordId
			}
			if count, ok := forwardMsg["count"].(float64); ok {
				rpcMsg.ForwardMsg.Count = int32(count)
			}
		}
	case 13:
		if markdownMsg, ok := msgData["markdownMsg"].(map[string]interface{}); ok {
			rpcMsg.MarkdownMsg = &chat_rpc.MarkdownMsg{}
			if content, ok := markdownMsg["content"].(string); ok {
				rpcMsg.MarkdownMsg.Content = content
			}
			if title, ok := markdownMsg["title"].(string); ok {
				rpcMsg.MarkdownMsg.Title = title
			}
		}
	case 14:
		if linkMsg, ok := msgData["linkMsg"].(map[string]interface{}); ok {
			rpcMsg.LinkMsg = &chat_rpc.LinkMsg{}
			if url, ok := linkMsg["url"].(string); ok {
				rpcMsg.LinkMsg.Url = url
			}
			if title, ok := linkMsg["title"].(string); ok {
				rpcMsg.LinkMsg.Title = title
			}
			if desc, ok := linkMsg["desc"].(string); ok {
				rpcMsg.LinkMsg.Desc = desc
			}
			if imageUrl, ok := linkMsg["imageUrl"].(string); ok {
				rpcMsg.LinkMsg.ImageUrl = imageUrl
			}
		}
	case 15:
		if cloudDocMsg, ok := msgData["cloudDocMsg"].(map[string]interface{}); ok {
			rpcMsg.CloudDocMsg = &chat_rpc.CloudDocMsg{}
			if docId, ok := cloudDocMsg["docId"].(string); ok {
				rpcMsg.CloudDocMsg.DocId = docId
			}
			if docType, ok := cloudDocMsg["docType"].(float64); ok {
				rpcMsg.CloudDocMsg.DocType = int32(docType)
			}
			if title, ok := cloudDocMsg["title"].(string); ok {
				rpcMsg.CloudDocMsg.Title = title
			}
			if ownerId, ok := cloudDocMsg["ownerId"].(string); ok {
				rpcMsg.CloudDocMsg.OwnerId = ownerId
			}
			if perm, ok := cloudDocMsg["perm"].(float64); ok {
				rpcMsg.CloudDocMsg.Perm = int32(perm)
			}
			if coverUrl, ok := cloudDocMsg["coverUrl"].(string); ok {
				rpcMsg.CloudDocMsg.CoverUrl = coverUrl
			}
			if revision, ok := cloudDocMsg["revision"].(float64); ok {
				rpcMsg.CloudDocMsg.Revision = int64(revision)
			}
		}
	case 16:
		if cardMsg, ok := msgData["cardMsg"].(map[string]interface{}); ok {
			rpcMsg.CardMsg = &chat_rpc.CardMsg{}
			if cardType, ok := cardMsg["cardType"].(float64); ok {
				rpcMsg.CardMsg.CardType = int32(cardType)
			}
			if id, ok := cardMsg["id"].(string); ok {
				rpcMsg.CardMsg.Id = id
			}
			if expireAt, ok := cardMsg["expireAt"].(float64); ok {
				rpcMsg.CardMsg.ExpireAt = int64(expireAt)
			}
			if inviteToken, ok := cardMsg["inviteToken"].(string); ok {
				rpcMsg.CardMsg.InviteToken = inviteToken
			}
		}
	}

	return rpcMsg, nil
}
