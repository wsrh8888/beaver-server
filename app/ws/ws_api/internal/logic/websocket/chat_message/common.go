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

			// 提取 fileKey（优先）或 fileName（兼容旧格式）
			if fileKey, ok := imageMsg["fileKey"].(string); ok {
				rpcMsg.ImageMsg.FileKey = fileKey
			} else if fileName, ok := imageMsg["fileName"].(string); ok {
				rpcMsg.ImageMsg.FileKey = fileName
			}

			// 提取打平后的字段（兼容 style 对象和直接字段两种格式）
			if style, ok := imageMsg["style"].(map[string]interface{}); ok {
				// 兼容旧格式：从 style 对象提取
				if width, ok := style["width"].(float64); ok {
					rpcMsg.ImageMsg.Width = int32(width)
				}
				if height, ok := style["height"].(float64); ok {
					rpcMsg.ImageMsg.Height = int32(height)
				}
			} else {
				// 新格式：直接从 imageMsg 提取 width 和 height
				if width, ok := imageMsg["width"].(float64); ok {
					rpcMsg.ImageMsg.Width = int32(width)
				}
				if height, ok := imageMsg["height"].(float64); ok {
					rpcMsg.ImageMsg.Height = int32(height)
				}
			}
			// 提取 size
			if size, ok := imageMsg["size"].(float64); ok {
				rpcMsg.ImageMsg.Size = int64(size)
			}
		}
	case 3: // 视频消息
		if videoMsg, ok := msgData["videoMsg"].(map[string]interface{}); ok {
			rpcMsg.VideoMsg = &chat_rpc.VideoMsg{}

			// 提取 fileKey（优先）或 fileName（兼容旧格式）
			if fileKey, ok := videoMsg["fileKey"].(string); ok {
				rpcMsg.VideoMsg.FileKey = fileKey
			} else if fileName, ok := videoMsg["fileName"].(string); ok {
				rpcMsg.VideoMsg.FileKey = fileName
			}

			// 提取打平后的字段（兼容 style 对象和直接字段两种格式）
			if style, ok := videoMsg["style"].(map[string]interface{}); ok {
				// 兼容旧格式：从 style 对象提取
				if width, ok := style["width"].(float64); ok {
					rpcMsg.VideoMsg.Width = int32(width)
				}
				if height, ok := style["height"].(float64); ok {
					rpcMsg.VideoMsg.Height = int32(height)
				}
				if duration, ok := style["duration"].(float64); ok {
					rpcMsg.VideoMsg.Duration = int32(duration)
				}
			} else {
				// 新格式：直接从 videoMsg 提取 width、height 和 duration
				if width, ok := videoMsg["width"].(float64); ok {
					rpcMsg.VideoMsg.Width = int32(width)
				}
				if height, ok := videoMsg["height"].(float64); ok {
					rpcMsg.VideoMsg.Height = int32(height)
				}
				if duration, ok := videoMsg["duration"].(float64); ok {
					rpcMsg.VideoMsg.Duration = int32(duration)
				}
			}
			// 提取 thumbnailKey 和 size
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
			// 提取 fileKey（优先）或 fileName（兼容旧格式）
			if fileKey, ok := fileMsg["fileKey"].(string); ok {
				file.FileKey = fileKey
			} else if fileName, ok := fileMsg["fileName"].(string); ok {
				file.FileKey = fileName
			}
			// 提取可选字段
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
	case 5: // 语音消息
		if voiceMsg, ok := msgData["voiceMsg"].(map[string]interface{}); ok {
			rpcMsg.VoiceMsg = &chat_rpc.VoiceMsg{}

			// 提取 fileKey（优先）或 fileName（兼容旧格式）
			if fileKey, ok := voiceMsg["fileKey"].(string); ok {
				rpcMsg.VoiceMsg.FileKey = fileKey
			} else if fileName, ok := voiceMsg["fileName"].(string); ok {
				rpcMsg.VoiceMsg.FileKey = fileName
			} else if src, ok := voiceMsg["src"].(string); ok {
				// 兼容旧格式：src 字段
				rpcMsg.VoiceMsg.FileKey = src
			}

			// 提取打平后的字段（兼容 style 对象和直接字段两种格式）
			if style, ok := voiceMsg["style"].(map[string]interface{}); ok {
				// 兼容旧格式：从 style 对象提取
				if duration, ok := style["duration"].(float64); ok {
					rpcMsg.VoiceMsg.Duration = int32(duration)
				}
			} else {
				// 新格式：直接从 voiceMsg 提取 duration
				if duration, ok := voiceMsg["duration"].(float64); ok {
					rpcMsg.VoiceMsg.Duration = int32(duration)
				}
			}
			// 提取 size
			if size, ok := voiceMsg["size"].(float64); ok {
				rpcMsg.VoiceMsg.Size = int64(size)
			}
		}
	case 6: // 表情消息
		if emojiMsg, ok := msgData["emojiMsg"].(map[string]interface{}); ok {
			emoji := &chat_rpc.EmojiMsg{}
			// 提取 fileKey（优先）或 fileName（兼容旧格式）
			if fileKey, ok := emojiMsg["fileKey"].(string); ok {
				emoji.FileKey = fileKey
			} else if fileName, ok := emojiMsg["fileName"].(string); ok {
				emoji.FileKey = fileName
			}
			if emojiId, ok := emojiMsg["emojiId"].(float64); ok {
				emoji.EmojiId = uint32(emojiId)
			}
			if packageId, ok := emojiMsg["packageId"].(float64); ok {
				emoji.PackageId = uint32(packageId)
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

			// 提取 fileKey（优先）或 fileName（兼容旧格式）
			if fileKey, ok := audioFileMsg["fileKey"].(string); ok {
				rpcMsg.AudioFileMsg.FileKey = fileKey
			} else if fileName, ok := audioFileMsg["fileName"].(string); ok {
				rpcMsg.AudioFileMsg.FileKey = fileName
			}

			// 提取可选字段
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
	}

	return rpcMsg, nil
}

// buildResponseData 构建响应数据
func buildResponseData(rpcResp *chat_rpc.SendMsgRes, originalMsg json.RawMessage) ([]byte, error) {
	responseData := map[string]interface{}{
		"id":             rpcResp.Id,
		"messageId":      rpcResp.MessageId,
		"conversationId": rpcResp.ConversationId,
		"msg":            originalMsg, // 使用原始消息数据
		"sender": map[string]interface{}{
			"userId":   rpcResp.Sender.UserId,
			"avatar":   rpcResp.Sender.Avatar,
			"nickname": rpcResp.Sender.Nickname,
		},
		"conversationType": rpcResp.ConversationType,
		"createAt":         rpcResp.CreateAt,
		"msgPreview":       rpcResp.MsgPreview,
		"status":           rpcResp.Status,
		"seq":              rpcResp.Seq,
	}

	// 将响应数据转换为 JSON 格式
	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		return nil, err
	}

	return responseJSON, nil
}
