package chat_message

import (
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"encoding/json"
)

// convertToRpcMsg 将原始消息转换为RPC消息格式
func convertToRpcMsg(msg json.RawMessage) (*chat_rpc.Msg, error) {
	var msgData map[string]interface{}
	err := json.Unmarshal(msg, &msgData)
	if err != nil {
		return nil, err
	}

	rpcMsg := &chat_rpc.Msg{}

	// 获取消息类型
	if msgType, ok := msgData["type"].(float64); ok {
		rpcMsg.Type = uint32(msgType)
	}

	// 根据消息类型设置对应的消息内容
	switch rpcMsg.Type {
	case 1: // 文本消息
		if textMsg, ok := msgData["textMsg"].(map[string]interface{}); ok {
			if content, ok := textMsg["content"].(string); ok {
				rpcMsg.TextMsg = &chat_rpc.TextMsg{Content: content}
			}
		}
	case 2: // 图片消息
		if imageMsg, ok := msgData["imageMsg"].(map[string]interface{}); ok {
			rpcMsg.ImageMsg = &chat_rpc.ImageMsg{}
			if fileKey, ok := imageMsg["fileKey"].(string); ok {
				rpcMsg.ImageMsg.FileKey = fileKey
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
	case 3: // 视频消息
		if videoMsg, ok := msgData["videoMsg"].(map[string]interface{}); ok {
			rpcMsg.VideoMsg = &chat_rpc.VideoMsg{}
			if fileKey, ok := videoMsg["fileKey"].(string); ok {
				rpcMsg.VideoMsg.FileKey = fileKey
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
			if thumbnailKey, ok := videoMsg["thumbnailKey"].(string); ok {
				rpcMsg.VideoMsg.ThumbnailKey = thumbnailKey
			}
			if size, ok := videoMsg["size"].(float64); ok {
				rpcMsg.VideoMsg.Size = int64(size)
			}
		}
	case 4: // 文件消息
		if fileMsg, ok := msgData["fileMsg"].(map[string]interface{}); ok {
			file := &chat_rpc.FileMsg{}
			if fileKey, ok := fileMsg["fileKey"].(string); ok {
				file.FileKey = fileKey
			}
			if size, ok := fileMsg["size"].(float64); ok {
				file.Size = int64(size)
			}
			if mimeType, ok := fileMsg["mimeType"].(string); ok {
				file.MimeType = mimeType
			}
			rpcMsg.FileMsg = file
		}
	case 5: // 语音消息
		if voiceMsg, ok := msgData["voiceMsg"].(map[string]interface{}); ok {
			rpcMsg.VoiceMsg = &chat_rpc.VoiceMsg{}
			if fileKey, ok := voiceMsg["fileKey"].(string); ok {
				rpcMsg.VoiceMsg.FileKey = fileKey
			}
			if duration, ok := voiceMsg["duration"].(float64); ok {
				rpcMsg.VoiceMsg.Duration = int32(duration)
			}
			if size, ok := voiceMsg["size"].(float64); ok {
				rpcMsg.VoiceMsg.Size = int64(size)
			}
		}
	case 6: // 表情消息
		if emojiMsg, ok := msgData["emojiMsg"].(map[string]interface{}); ok {
			emoji := &chat_rpc.EmojiMsg{}
			if fileKey, ok := emojiMsg["fileKey"].(string); ok {
				emoji.FileKey = fileKey
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
	case 7: // 通知消息
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
	case 8: // 音频文件消息
		if audioFileMsg, ok := msgData["audioFileMsg"].(map[string]interface{}); ok {
			rpcMsg.AudioFileMsg = &chat_rpc.AudioFileMsg{}
			if fileKey, ok := audioFileMsg["fileKey"].(string); ok {
				rpcMsg.AudioFileMsg.FileKey = fileKey
			}
			if duration, ok := audioFileMsg["duration"].(float64); ok {
				rpcMsg.AudioFileMsg.Duration = int32(duration)
			}
			if size, ok := audioFileMsg["size"].(float64); ok {
				rpcMsg.AudioFileMsg.Size = int64(size)
			}
		}
	case 9: // 音视频通话
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
	case 10: // 撤回消息
		if withdrawMsg, ok := msgData["withdrawMsg"].(map[string]interface{}); ok {
			rpcMsg.WithdrawMsg = &chat_rpc.WithdrawMsg{}
			if originMsgId, ok := withdrawMsg["originMsgId"].(string); ok {
				rpcMsg.WithdrawMsg.OriginMsgId = originMsgId
			}
			if originMsgMap, ok := withdrawMsg["originMsg"].(map[string]interface{}); ok {
				// 递归转换快照内容
				jsonData, _ := json.Marshal(originMsgMap)
				originMsg, _ := convertToRpcMsg(jsonData)
				rpcMsg.WithdrawMsg.OriginMsg = originMsg
			}
		}
	case 11: // 回复消息
		if replyMsg, ok := msgData["replyMsg"].(map[string]interface{}); ok {
			rpcMsg.ReplyMsg = &chat_rpc.ReplyMsg{}
			if originMsgId, ok := replyMsg["originMsgId"].(string); ok {
				rpcMsg.ReplyMsg.OriginMsgId = originMsgId
			}
			if originMsgMap, ok := replyMsg["originMsg"].(map[string]interface{}); ok {
				// 递归转换被引用的消息快照
				jsonData, _ := json.Marshal(originMsgMap)
				originMsg, _ := convertToRpcMsg(jsonData)
				rpcMsg.ReplyMsg.OriginMsg = originMsg
			}
			if replyMsgInnerMap, ok := replyMsg["replyMsg"].(map[string]interface{}); ok {
				// 递归转换回复的具体消息内容
				jsonData, _ := json.Marshal(replyMsgInnerMap)
				replyInnerMsg, _ := convertToRpcMsg(jsonData)
				rpcMsg.ReplyMsg.ReplyMsg = replyInnerMsg
			}
		}
	case 12: // 转发消息
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
	}

	return rpcMsg, nil
}
