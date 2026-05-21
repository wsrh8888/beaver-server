package robot

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type RobotSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群自定义机器人消息推送（Jenkins/GitHub/Grafana 等）
func NewRobotSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RobotSendLogic {
	return &RobotSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RobotSendLogic) RobotSend(req *types.RobotSendReq) (resp *types.RobotSendRes, err error) {
	if req.AccessToken == "" {
		return nil, errors.New("无效的 access_token")
	}

	var webhook open_models.OpenGroupBotModel
	if err := l.svcCtx.DB.Where("token = ? AND status = ?", req.AccessToken, 1).First(&webhook).Error; err != nil {
		return nil, errors.New("无效的 access_token")
	}

	nowMs := time.Now().UnixMilli()
	if req.Timestamp <= 0 || robotAbs64(nowMs-req.Timestamp) > 3600000 {
		return nil, errors.New("请求已过期")
	}

	sign, decodeErr := url.QueryUnescape(req.Sign)
	if decodeErr != nil {
		sign = req.Sign
	}
	if !verifyRobotSign(req.Timestamp, webhook.Secret, sign) {
		return nil, errors.New("签名校验失败")
	}

	rpcMsg, buildErr := l.buildRpcMsg(req)
	if buildErr != nil {
		return nil, buildErr
	}

	sendReq := &chat_rpc.SendMsgReq{
		UserId:         webhook.BotUserID,
		ConversationId: webhook.GroupID,
		Msg:            rpcMsg,
	}
	if _, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, sendReq); err != nil {
		logx.Errorf("Robot 消息发送失败: group=%s bot=%s err=%v", webhook.GroupID, webhook.BotUserID, err)
		return nil, errors.New("消息发送失败")
	}

	logx.Infof("Robot 推送成功: group=%s msgtype=%s", webhook.GroupID, req.Msgtype)
	return &types.RobotSendRes{}, nil
}

func verifyRobotSign(timestamp int64, secret, sign string) bool {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(stringToSign))
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(sign))
}

func (l *RobotSendLogic) buildRpcMsg(req *types.RobotSendReq) (*chat_rpc.Msg, error) {
	atSuffix, atUserIDs := l.resolveAtUsers(req.At)
	msg, err := l.buildRpcMsgBody(req, atSuffix)
	if err != nil {
		return nil, err
	}
	if len(atUserIDs) > 0 {
		msg.AtUserIds = atUserIDs
	}
	return msg, nil
}

func (l *RobotSendLogic) buildRpcMsgBody(req *types.RobotSendReq, atSuffix string) (*chat_rpc.Msg, error) {
	switch req.Msgtype {
	// ── 1. 文本消息 ──────────────────────────────────────────────────────
	case "text":
		content := strings.TrimSpace(req.Text.Content)
		if content == "" {
			return nil, errors.New("text.content 不能为空")
		}
		return &chat_rpc.Msg{
			Type:    1,
			TextMsg: &chat_rpc.TextMsg{Content: robotAppendAt(content, atSuffix)},
		}, nil

	// ── 2. Markdown 消息 ─────────────────────────────────────────────────
	case "markdown":
		title := strings.TrimSpace(req.Markdown.Title)
		text := strings.TrimSpace(req.Markdown.Text)
		if title == "" || text == "" {
			return nil, errors.New("markdown.title 与 markdown.text 不能为空")
		}
		var sb strings.Builder
		sb.WriteString(text)
		if img := strings.TrimSpace(req.Markdown.Image); img != "" {
			sb.WriteString("\n\n![image](")
			sb.WriteString(img)
			sb.WriteString(")")
		}
		if atSuffix != "" {
			sb.WriteString(atSuffix)
		}
		return &chat_rpc.Msg{
			Type:        13,
			MarkdownMsg: &chat_rpc.MarkdownMsg{Title: title, Content: sb.String()},
		}, nil

	// ── 3. 图片消息 ──────────────────────────────────────────────────────
	case "image":
		imgURL := strings.TrimSpace(req.Image.Url)
		if imgURL == "" {
			return nil, errors.New("image.url 不能为空")
		}
		width, _ := strconv.Atoi(req.Image.Width)
		height, _ := strconv.Atoi(req.Image.Height)
		return &chat_rpc.Msg{
			Type: 2,
			ImageMsg: &chat_rpc.ImageMsg{
				FileKey: imgURL,
				Width:   int32(width),
				Height:  int32(height),
			},
		}, nil

	// ── 4. 链接消息 ──────────────────────────────────────────────────────
	case "link":
		msgURL := strings.TrimSpace(req.Link.MessageUrl)
		title := strings.TrimSpace(req.Link.Title)
		text := strings.TrimSpace(req.Link.Text)
		if msgURL == "" || title == "" || text == "" {
			return nil, errors.New("link.message_url, link.title, link.text 不能为空")
		}
		return &chat_rpc.Msg{
			Type: 14,
			LinkMsg: &chat_rpc.LinkMsg{
				Url:      msgURL,
				Title:    title,
				Desc:     text,
				ImageUrl: strings.TrimSpace(req.Link.PicUrl),
			},
		}, nil

	// ── 5. 语音消息 ──────────────────────────────────────────────────────
	case "audio":
		audioURL := strings.TrimSpace(req.Audio.Url)
		if audioURL == "" {
			return nil, errors.New("audio.url 不能为空")
		}
		if req.Audio.Duration <= 0 {
			return nil, errors.New("audio.duration 必须为正整数（秒）")
		}
		return &chat_rpc.Msg{
			Type: 5,
			VoiceMsg: &chat_rpc.VoiceMsg{
				FileKey:  audioURL,
				Duration: int32(req.Audio.Duration),
				Size:     req.Audio.Size * 1024, // KB → bytes
			},
		}, nil

	// ── 6. 文件消息 ──────────────────────────────────────────────────────
	case "file":
		fileURL := strings.TrimSpace(req.File.Url)
		fileName := strings.TrimSpace(req.File.Name)
		if fileURL == "" || fileName == "" {
			return nil, errors.New("file.url, file.name 不能为空")
		}
		return &chat_rpc.Msg{
			Type: 4,
			FileMsg: &chat_rpc.FileMsg{
				FileKey:  fileURL,
				FileName: fileName,
				Size:     req.File.Size * 1024, // KB → bytes
			},
		}, nil

	// ── 7. 视频消息 ──────────────────────────────────────────────────────
	case "video":
		videoURL := strings.TrimSpace(req.Video.Url)
		videoName := strings.TrimSpace(req.Video.Name)
		if videoURL == "" || videoName == "" {
			return nil, errors.New("video.url, video.name 不能为空")
		}
		width, _ := strconv.Atoi(req.Video.Width)
		height, _ := strconv.Atoi(req.Video.Height)
		sizeKB, _ := strconv.ParseInt(strings.TrimSpace(req.Video.Size), 10, 64)
		return &chat_rpc.Msg{
			Type: 3,
			VideoMsg: &chat_rpc.VideoMsg{
				FileKey:  videoURL,
				Duration: int32(req.Video.Duration),
				Width:    int32(width),
				Height:   int32(height),
				Size:     sizeKB * 1024, // KB → bytes
			},
		}, nil

	// ── 8. 页面跳转消息 ──────────────────────────────────────────────────
	case "custom":
		bodyURL := strings.TrimSpace(req.Custom.BodyUrl)
		if bodyURL == "" {
			return nil, errors.New("custom.body_url 不能为空")
		}
		return &chat_rpc.Msg{
			Type:    14,
			LinkMsg: &chat_rpc.LinkMsg{Url: bodyURL, Title: "打开页面"},
		}, nil

	// ── 9. 灰色 Tips 消息 ────────────────────────────────────────────────
	case "tips":
		text := strings.TrimSpace(req.Tips.Text)
		if text == "" {
			return nil, errors.New("tips.text 不能为空")
		}
		return &chat_rpc.Msg{
			Type:    1,
			TextMsg: &chat_rpc.TextMsg{Content: text},
		}, nil

	// ── 10. ActionCard 卡片消息 ───────────────────────────────────────────
	case "action_card":
		title := strings.TrimSpace(req.ActionCard.Title)
		markdown := strings.TrimSpace(req.ActionCard.Markdown)
		if title == "" || markdown == "" {
			return nil, errors.New("action_card.title, action_card.markdown 不能为空")
		}
		var sb strings.Builder
		sb.WriteString("## ")
		sb.WriteString(title)
		sb.WriteString("\n")
		if ct := strings.TrimSpace(req.ActionCard.ContentTitle); ct != "" {
			sb.WriteString("**")
			sb.WriteString(ct)
			sb.WriteString("**\n")
		}
		sb.WriteString(markdown)
		if img := strings.TrimSpace(req.ActionCard.Image); img != "" {
			sb.WriteString("\n\n![](")
			sb.WriteString(img)
			sb.WriteString(")")
		}
		if su := strings.TrimSpace(req.ActionCard.SingleUrl); su != "" {
			btnText := req.ActionCard.SingleTitle
			if strings.TrimSpace(btnText) == "" {
				btnText = "查看详情"
			}
			sb.WriteString("\n\n[")
			sb.WriteString(btnText)
			sb.WriteString("](")
			sb.WriteString(su)
			sb.WriteString(")")
		} else if len(req.ActionCard.BtnJsonList) > 0 {
			sb.WriteString("\n")
			for _, btn := range req.ActionCard.BtnJsonList {
				btnTitle := strings.TrimSpace(btn.Title)
				btnURL := strings.TrimSpace(btn.ActionUrl)
				if btnTitle == "" || btnURL == "" {
					continue
				}
				sb.WriteString("\n[")
				sb.WriteString(btnTitle)
				sb.WriteString("](")
				sb.WriteString(btnURL)
				sb.WriteString(")")
			}
		}
		return &chat_rpc.Msg{
			Type:        13,
			MarkdownMsg: &chat_rpc.MarkdownMsg{Title: title, Content: robotAppendAt(sb.String(), atSuffix)},
		}, nil

	default:
		return nil, fmt.Errorf("不支持的消息类型: %s，支持 text/markdown/image/link/audio/file/video/custom/tips/action_card", req.Msgtype)
	}
}

func (l *RobotSendLogic) resolveAtUsers(at types.RobotSendAt) (suffix string, userIDs []string) {
	if at.IsAtAll {
		return "\n@所有人", nil
	}

	seen := make(map[string]struct{})
	addID := func(id string) {
		id = strings.TrimSpace(id)
		if id == "" {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		userIDs = append(userIDs, id)
	}

	for _, uid := range at.AtUserIds {
		addID(uid)
	}

	for _, wc := range at.AtWorkCodes {
		addID(wc)
	}

	for _, mobile := range at.AtMobiles {
		mobile = strings.TrimSpace(mobile)
		if mobile == "" {
			continue
		}
		var user user_models.UserModel
		if err := l.svcCtx.DB.Where("phone = ? AND status = 1", mobile).First(&user).Error; err == nil {
			addID(user.UserID)
		}
	}

	if len(userIDs) == 0 {
		return "", nil
	}
	var sb strings.Builder
	sb.WriteString("\n")
	for _, uid := range userIDs {
		sb.WriteString("@")
		sb.WriteString(uid)
		sb.WriteString(" ")
	}
	return sb.String(), userIDs
}

func robotAppendAt(content, atSuffix string) string {
	if atSuffix == "" {
		return content
	}
	return content + atSuffix
}

func robotAbs64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
