package bot_public

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
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

func (l *BotSendLogic) BotSend(req *types.BotSendReq, clientIP string) (resp *types.BotSendRes, err error) {
	// 1. 根据 Token 查询机器人信息
	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("token = ?", req.Token).First(&bot).Error; err != nil {
		l.Errorf("BotSend: query bot failed, token=%s, error=%v", req.Token, err)
		return nil, fmt.Errorf("invalid token")
	}

	l.Infof("BotSend: found bot, botID=%s, groupID=%s, status=%d", bot.BotID, bot.GroupID, bot.Status)

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

		// 验证签名：HMAC-SHA256 + Base64
		expectedSign := generateSignature(req.Timestamp, bot.Security.SignatureSecret)
		if req.Sign != expectedSign {
			l.Errorf("BotSend: signature mismatch, received=%s, expected=%s", req.Sign, expectedSign)
			return nil, fmt.Errorf("invalid signature")
		}
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
		if clientIP == "" {
			return nil, fmt.Errorf("client ip required")
		}
		host := strings.TrimSpace(clientIP)
		if h, _, splitErr := net.SplitHostPort(host); splitErr == nil {
			host = h
		}
		allowed := false
		for _, item := range bot.Security.IPWhitelist {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			if item == host || item == clientIP {
				allowed = true
				break
			}
		}
		if !allowed {
			l.Errorf("BotSend: ip not in whitelist, client=%s botID=%s", host, bot.BotID)
			return nil, fmt.Errorf("ip not in whitelist")
		}
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
		msg.Type = 13 // Markdown 消息
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
			FileUrl: req.Image.URL,
			Width:   int32(req.Image.Width),
			Height:  int32(req.Image.Height),
		}

	case "video":
		if req.Video == nil {
			return nil, fmt.Errorf("video content is required")
		}
		msg.Type = 3 // 视频消息
		msg.VideoMsg = &chat_rpc.VideoMsg{
			FileUrl:       req.Video.URL,
			Width:        int32(req.Video.Width),
			Height:       int32(req.Video.Height),
			Duration:     int32(req.Video.Duration),
			ThumbnailUrl: req.Video.ThumbnailURL,
		}

	case "file":
		if req.File == nil {
			return nil, fmt.Errorf("file content is required")
		}
		msg.Type = 4 // 文件消息
		msg.FileMsg = &chat_rpc.FileMsg{
			FileUrl:  req.File.URL,
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
			FileUrl:  req.Voice.URL,
			Duration: int32(req.Voice.Duration),
			Size:     req.Voice.FileSize,
		}

	case "audio_file":
		if req.AudioFile == nil {
			return nil, fmt.Errorf("audio file content is required")
		}
		msg.Type = 8 // 音频文件消息
		msg.AudioFileMsg = &chat_rpc.AudioFileMsg{
			FileUrl:  req.AudioFile.URL,
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
			FileUrl: req.Emoji.URL,
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
		if req.Reply.OriginMsgID == "" {
			return nil, fmt.Errorf("reply originMsgId is required")
		}
		replyBody := &chat_rpc.Msg{}
		switch req.Reply.MsgType {
		case "text", "":
			replyBody.Type = 1
			replyBody.TextMsg = &chat_rpc.TextMsg{Content: req.Reply.Content}
		case "markdown":
			replyBody.Type = 13
			replyBody.MarkdownMsg = &chat_rpc.MarkdownMsg{Content: req.Reply.Content}
		default:
			return nil, fmt.Errorf("unsupported reply msgType: %s", req.Reply.MsgType)
		}
		originSnap := &chat_rpc.Msg{
			Type:    1,
			TextMsg: &chat_rpc.TextMsg{Content: "[消息]"},
		}
		convID := "group_" + bot.GroupID
		getRes, getErr := l.svcCtx.ChatRpc.GetChatMessage(l.ctx, &chat_rpc.GetChatMessageReq{
			ConversationId: convID,
			MessageId:      req.Reply.OriginMsgID,
		})
		if getErr == nil && getRes != nil && getRes.Found && getRes.Msg != nil {
			originSnap = getRes.Msg
		}
		msg.Type = 11
		msg.ReplyMsg = &chat_rpc.ReplyMsg{
			OriginMsgId: req.Reply.OriginMsgID,
			OriginMsg:   originSnap,
			ReplyMsg:    replyBody,
		}

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

	default:
		return nil, fmt.Errorf("unsupported message type: %s", req.MsgType)
	}

	// 设置 @ 用户列表
	msg.AtUserIds = atUserIds

	// 7. 生成消息 ID
	messageId := fmt.Sprintf("bot_%d_%d", bot.ID, time.Now().UnixNano())

	// 8. 调用 chat_rpc 发送消息
	// ConversationId 增加 group_ 前缀表示群聊
	chatReq := &chat_rpc.SendMsgReq{
		UserId:         bot.BotID, // 使用机器人的 UserID 作为发送者
		ConversationId: "group_" + bot.GroupID,
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

// generateSignature 生成签名（HMAC-SHA256 + Base64）
// 签名字符串 = timestamp + "\n" + secret
func generateSignature(timestamp int64, secret string) string {

	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(stringToSign))
	signature := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}
