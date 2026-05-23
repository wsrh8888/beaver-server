package bot

import (
	"context"
	"fmt"
	"time"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type BotSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 推送机器人发送消息到群（第三方服务如 Jenkins/GitLab 调用此接口）
func NewBotSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BotSendLogic {
	return &BotSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BotSendLogic) BotSend(req *types.BotSendReq) (resp *types.BotSendRes, err error) {
	// 1. 根据 Token 查询机器人信息
	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("token = ?", req.Token).First(&bot).Error; err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	// 2. 校验机器人状态
	if bot.Status != 1 {
		return nil, fmt.Errorf("bot is disabled")
	}

	// 3. 安全校验（如果启用了签名校验）
	if bot.Security.SignatureEnabled {
		if req.Timestamp == 0 || req.Sign == "" {
			return nil, fmt.Errorf("timestamp and sign are required")
		}

		// 校验时间戳（1小时有效期）
		now := time.Now().UnixMilli()
		if now-req.Timestamp > 3600000 {
			return nil, fmt.Errorf("timestamp expired")
		}

		// TODO: 验证签名（需要根据实际签名算法实现）
		// 这里简化处理，实际应该计算 HMAC-SHA256 签名并比对
	}

	// 4. 校验关键词（如果启用了关键词校验）
	if bot.Security.KeywordsEnabled && len(bot.Security.Keywords) > 0 {
		// 获取消息内容用于关键词匹配
		content := ""
		switch req.MsgType {
		case "text":
			if req.Text != nil {
				content = req.Text.Content
			}
		case "markdown":
			if req.Markdown != nil {
				content = req.Markdown.Content
			}
		}

		// 检查是否包含任一关键词
		hasKeyword := false
		for _, keyword := range bot.Security.Keywords {
			if content != "" && contains(content, keyword) {
				hasKeyword = true
				break
			}
		}
		if !hasKeyword {
			return nil, fmt.Errorf("message does not contain required keywords")
		}
	}

	// 5. 校验 IP 白名单（如果启用了 IP 白名单）
	if bot.Security.IPWhitelistEnabled && len(bot.Security.IPWhitelist) > 0 {
		// TODO: 获取请求 IP 并校验是否在白名单中
		// 这里需要通过 context 或 middleware 传递请求 IP
	}

	// 6. 构建 chat_rpc 的消息对象
	msg := &chat_rpc.Msg{}
	atUserIds := []string{}

	// 根据消息类型构建不同的消息内容
	switch req.MsgType {
	case "text":
		if req.Text == nil {
			return nil, fmt.Errorf("text content is required")
		}
		msg.Type = 1 // 文本消息
		msg.TextMsg = &chat_rpc.TextMsg{
			Content: req.Text.Content,
		}
		if req.Text.At != nil {
			atUserIds = req.Text.At.AtUserIDs
		}

	case "markdown":
		if req.Markdown == nil {
			return nil, fmt.Errorf("markdown content is required")
		}
		msg.Type = 15 // Markdown 消息
		msg.MarkdownMsg = &chat_rpc.MarkdownMsg{
			Title:   req.Markdown.Title,
			Content: req.Markdown.Content,
		}
		if req.Markdown.At != nil {
			atUserIds = req.Markdown.At.AtUserIDs
		}

	case "image":
		if req.Image == nil {
			return nil, fmt.Errorf("image content is required")
		}
		msg.Type = 2 // 图片消息
		msg.ImageMsg = &chat_rpc.ImageMsg{
			FileKey: req.Image.URL,
			Width:   int32(req.Image.Width),
			Height:  int32(req.Image.Height),
		}

	case "video":
		if req.Video == nil {
			return nil, fmt.Errorf("video content is required")
		}
		msg.Type = 3 // 视频消息
		msg.VideoMsg = &chat_rpc.VideoMsg{
			FileKey:      req.Video.URL,
			Width:        int32(req.Video.Width),
			Height:       int32(req.Video.Height),
			Duration:     int32(req.Video.Duration),
			ThumbnailKey: req.Video.ThumbnailURL,
		}

	case "file":
		if req.File == nil {
			return nil, fmt.Errorf("file content is required")
		}
		msg.Type = 4 // 文件消息
		msg.FileMsg = &chat_rpc.FileMsg{
			FileKey:  req.File.URL,
			FileName: req.File.FileName,
			Size:     req.File.FileSize,
			MimeType: req.File.MimeType,
		}

	case "voice":
		if req.Voice == nil {
			return nil, fmt.Errorf("voice content is required")
		}
		msg.Type = 5 // 语音消息
		msg.VoiceMsg = &chat_rpc.VoiceMsg{
			FileKey:  req.Voice.URL,
			Duration: int32(req.Voice.Duration),
			Size:     req.Voice.FileSize,
		}

	case "audio_file":
		if req.AudioFile == nil {
			return nil, fmt.Errorf("audio file content is required")
		}
		msg.Type = 8 // 音频文件消息
		msg.AudioFileMsg = &chat_rpc.AudioFileMsg{
			FileKey:  req.AudioFile.URL,
			FileName: req.AudioFile.FileName,
			Duration: int32(req.AudioFile.Duration),
			Size:     req.AudioFile.FileSize,
		}

	case "emoji":
		if req.Emoji == nil {
			return nil, fmt.Errorf("emoji content is required")
		}
		msg.Type = 6 // 表情消息
		msg.EmojiMsg = &chat_rpc.EmojiMsg{
			FileKey: req.Emoji.URL,
			Width:   req.Emoji.Width,
			Height:  req.Emoji.Height,
		}

	case "notification":
		if req.Notification == nil {
			return nil, fmt.Errorf("notification content is required")
		}
		msg.Type = 7 // 通知消息
		msg.NotificationMsg = &chat_rpc.NotificationMsg{
			Type:   req.Notification.Type,
			Actors: req.Notification.Actors,
		}

	case "call":
		if req.Call == nil {
			return nil, fmt.Errorf("call content is required")
		}
		msg.Type = 9 // 通话消息
		msg.CallMsg = &chat_rpc.CallMsg{
			RoomId:   req.Call.RoomID,
			CallType: req.Call.CallType,
			Status:   req.Call.Status,
			Duration: req.Call.Duration,
		}

	case "withdraw":
		if req.Withdraw == nil {
			return nil, fmt.Errorf("withdraw content is required")
		}
		msg.Type = 10 // 撤回消息
		msg.TargetMsgId = req.Withdraw.OriginMsgID

	case "reply":
		if req.Reply == nil {
			return nil, fmt.Errorf("reply content is required")
		}
		msg.Type = 11 // 回复消息
		// TODO: 需要查询原消息内容来构建 ReplyMsg

	case "forward":
		if req.Forward == nil {
			return nil, fmt.Errorf("forward content is required")
		}
		msg.Type = 12 // 转发消息
		msg.ForwardMsg = &chat_rpc.ForwardMsg{
			Title:    req.Forward.Title,
			RecordId: req.Forward.RecordID,
			Count:    req.Forward.Count,
		}

	case "link":
		if req.Link == nil {
			return nil, fmt.Errorf("link content is required")
		}
		msg.Type = 16 // 链接卡片消息
		msg.LinkMsg = &chat_rpc.LinkMsg{
			Url:      req.Link.URL,
			Title:    req.Link.Title,
			Desc:     req.Link.Desc,
			ImageUrl: req.Link.ImageURL,
		}

	case "card":
		if req.Card == nil {
			return nil, fmt.Errorf("card content is required")
		}
		// TODO: 交互式卡片消息需要特殊处理
		return nil, fmt.Errorf("card message type not fully supported yet")

	default:
		return nil, fmt.Errorf("unsupported message type: %s", req.MsgType)
	}

	// 设置 @ 用户列表
	msg.AtUserIds = atUserIds

	// 7. 生成消息 ID
	messageId := fmt.Sprintf("bot_%d_%d", bot.ID, time.Now().UnixNano())

	// 8. 调用 chat_rpc 发送消息
	chatReq := &chat_rpc.SendMsgReq{
		UserId:         bot.BotID, // 使用机器人的 UserID 作为发送者
		ConversationId: bot.GroupID,
		MessageId:      messageId,
		Msg:            msg,
		DeviceId:       "webhook",
	}

	chatRes, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, chatReq)
	if err != nil {
		l.Errorf("failed to send message via chat_rpc: %v", err)
		return nil, fmt.Errorf("failed to send message")
	}

	// 9. 返回响应
	return &types.BotSendRes{
		MessageID: chatRes.MessageId,
		SendTime:  parseTime(chatRes.CreatedAt),
	}, nil
}

// contains 检查字符串是否包含子串（简单实现）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// parseTime 解析时间字符串为毫秒时间戳
func parseTime(timeStr string) int64 {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return time.Now().UnixMilli()
	}
	return t.UnixMilli()
}
