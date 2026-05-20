package webhook

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

const (
	webhookCodeOK           = 0
	webhookCodeInvalidToken = 401
	webhookCodeExpired      = 408
	webhookCodeSignFailed   = 403
)

type IncomingWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIncomingWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IncomingWebhookLogic {
	return &IncomingWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IncomingWebhookLogic) IncomingWebhook(req *types.IncomingWebhookReq) (resp *types.IncomingWebhookRes, err error) {
	if req.AccessToken == "" {
		return webhookResp(webhookCodeInvalidToken, "无效的 access_token"), nil
	}

	var webhook open_models.OpenIncomingWebhook
	if err := l.svcCtx.DB.Where("token = ? AND status = ?", req.AccessToken, 1).First(&webhook).Error; err != nil {
		return webhookResp(webhookCodeInvalidToken, "无效的 access_token"), nil
	}

	nowMs := time.Now().UnixMilli()
	if req.Timestamp <= 0 || abs64(nowMs-req.Timestamp) > 3600000 {
		return webhookResp(webhookCodeExpired, "请求已过期"), nil
	}

	sign, decodeErr := url.QueryUnescape(req.Sign)
	if decodeErr != nil {
		sign = req.Sign
	}
	if !verifyWebhookSign(req.Timestamp, webhook.Secret, sign) {
		return webhookResp(webhookCodeSignFailed, "签名校验失败"), nil
	}

	rpcMsg, buildErr := l.buildRpcMsg(req)
	if buildErr != nil {
		return webhookResp(webhookCodeSignFailed, buildErr.Error()), nil
	}

	sendReq := &chat_rpc.SendMsgReq{
		UserId:         webhook.BotUserID,
		ConversationId: webhook.GroupID,
		Msg:            rpcMsg,
	}
	if _, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, sendReq); err != nil {
		logx.Errorf("Webhook 消息发送失败: group=%s bot=%s err=%v", webhook.GroupID, webhook.BotUserID, err)
		return webhookResp(webhookCodeSignFailed, "消息发送失败"), nil
	}

	logx.Infof("Webhook 推送成功: name=%s group=%s msgtype=%s", webhook.Name, webhook.GroupID, req.Msgtype)
	return webhookResp(webhookCodeOK, "ok"), nil
}

func webhookResp(code int, msg string) *types.IncomingWebhookRes {
	return &types.IncomingWebhookRes{Code: code, Msg: msg}
}

func verifyWebhookSign(timestamp int64, secret, sign string) bool {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(stringToSign))
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(sign))
}

func (l *IncomingWebhookLogic) buildRpcMsg(req *types.IncomingWebhookReq) (*chat_rpc.Msg, error) {
	atSuffix, _ := l.resolveAtUsers(req.At)

	switch req.Msgtype {
	case "text":
		content := strings.TrimSpace(req.Text.Content)
		if content == "" {
			return nil, errors.New("text.content 不能为空")
		}
		content = appendAtSuffix(content, atSuffix)
		return &chat_rpc.Msg{
			Type: 1,
			TextMsg: &chat_rpc.TextMsg{
				Content: content,
			},
		}, nil

	case "markdown":
		title := strings.TrimSpace(req.Markdown.Title)
		text := strings.TrimSpace(req.Markdown.Text)
		if title == "" || text == "" {
			return nil, errors.New("markdown.title 与 markdown.text 不能为空")
		}
		var sb strings.Builder
		sb.WriteString("## ")
		sb.WriteString(title)
		sb.WriteString("\n")
		sb.WriteString(text)
		if img := strings.TrimSpace(req.Markdown.Image); img != "" {
			sb.WriteString("\n\n![image](")
			sb.WriteString(img)
			sb.WriteString(")")
		}
		content := appendAtSuffix(sb.String(), atSuffix)
		return &chat_rpc.Msg{
			Type: 1,
			TextMsg: &chat_rpc.TextMsg{
				Content: content,
			},
		}, nil

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

	default:
		return nil, fmt.Errorf("不支持的消息类型: %s，仅支持 text / markdown / image", req.Msgtype)
	}
}

func (l *IncomingWebhookLogic) resolveAtUsers(at types.IncomingWebhookAt) (suffix string, userIDs []string) {
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

func appendAtSuffix(content, atSuffix string) string {
	if atSuffix == "" {
		return content
	}
	return content + atSuffix
}

func abs64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
