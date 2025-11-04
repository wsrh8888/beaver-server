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

			if fileName, ok := imageMsg["fileName"].(string); ok {
				rpcMsg.ImageMsg.FileKey = fileName
			}
			// 提取宽度和高度
			if width, ok := imageMsg["width"].(float64); ok {
				rpcMsg.ImageMsg.Width = int32(width)
			}
			if height, ok := imageMsg["height"].(float64); ok {
				rpcMsg.ImageMsg.Height = int32(height)
			}
		}
	case 3: // 视频消息
		if videoMsg, ok := msgData["videoMsg"].(map[string]interface{}); ok {
			rpcMsg.VideoMsg = &chat_rpc.VideoMsg{}

			if fileName, ok := videoMsg["fileName"].(string); ok {
				rpcMsg.VideoMsg.FileKey = fileName
			}

			// 提取宽度、高度和时长
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
	case 4: // 文件消息
		if fileMsg, ok := msgData["fileMsg"].(map[string]interface{}); ok {
			if fileName, ok := fileMsg["fileName"].(string); ok {
				rpcMsg.FileMsg = &chat_rpc.FileMsg{FileKey: fileName}
			}
		}
	case 5: // 语音消息
		if voiceMsg, ok := msgData["voiceMsg"].(map[string]interface{}); ok {
			rpcMsg.VoiceMsg = &chat_rpc.VoiceMsg{}

			if fileName, ok := voiceMsg["fileName"].(string); ok {
				rpcMsg.VoiceMsg.FileKey = fileName
			}

			// 提取时长
			if duration, ok := voiceMsg["duration"].(float64); ok {
				rpcMsg.VoiceMsg.Duration = int32(duration)
			}
		}
	case 6: // 表情消息
		if emojiMsg, ok := msgData["emojiMsg"].(map[string]interface{}); ok {
			emoji := &chat_rpc.EmojiMsg{}
			if fileName, ok := emojiMsg["fileName"].(string); ok {
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
