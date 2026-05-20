package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
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
	fail := func(msg string) (*types.IncomingWebhookRes, error) {
		return &types.IncomingWebhookRes{Success: false, ErrMsg: msg}, nil
	}

	// 1. 根据 Token 查询 Webhook 配置
	var webhook open_models.OpenIncomingWebhook
	if err := l.svcCtx.DB.Where("token = ? AND status = ?", req.Token, 1).First(&webhook).Error; err != nil {
		return fail("无效的 Webhook Token")
	}

	// 2. 防重放：时间戳必须在当前时间 ±5 分钟内
	now := time.Now().Unix()
	if req.Timestamp <= 0 || abs(now-req.Timestamp) > 300 {
		return fail("请求已过期，请检查服务器时间")
	}

	// 3. HMAC-SHA256 签名验证（与钉钉自定义机器人签名逻辑相同）
	//    stringToSign = timestamp + "\n" + secret
	//    sign = Base64( HMAC-SHA256( stringToSign, secret ) )
	if !verifySign(req.Timestamp, webhook.Secret, req.Sign) {
		return fail("签名验证失败")
	}

	// 4. 构造消息内容
	content, msgErr := buildContent(req)
	if msgErr != nil {
		return fail(msgErr.Error())
	}

	// 5. 通过 Chat RPC 发送消息到群
	sendReq := &chat_rpc.SendMsgReq{
		UserId:         webhook.BotUserID,
		ConversationId: webhook.GroupID,
		Msg: &chat_rpc.Msg{
			Type: 1,
			TextMsg: &chat_rpc.TextMsg{
				Content: content,
			},
		},
	}
	if _, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, sendReq); err != nil {
		logx.Errorf("Webhook 消息发送失败: group=%s, err=%v", webhook.GroupID, err)
		return fail("消息发送失败")
	}

	logx.Infof("Webhook 推送成功: name=%s group=%s", webhook.Name, webhook.GroupID)
	return &types.IncomingWebhookRes{Success: true}, nil
}

// verifySign 验证 HMAC-SHA256 签名
func verifySign(timestamp int64, secret, sign string) bool {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(stringToSign))
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(sign))
}

// buildContent 根据消息类型构造最终发送的文本内容
func buildContent(req *types.IncomingWebhookReq) (string, error) {
	if req.Content == "" {
		return "", errors.New("content 不能为空")
	}

	var sb strings.Builder

	switch req.MsgType {
	case "text":
		sb.WriteString(req.Content)
	case "markdown":
		if req.Title != "" {
			sb.WriteString("## ")
			sb.WriteString(req.Title)
			sb.WriteString("\n")
		}
		sb.WriteString(req.Content)
	default:
		return "", errors.New("不支持的消息类型，仅支持 text / markdown")
	}

	// 追加 @ 用户（如果有）
	if req.AtAll {
		sb.WriteString("\n@所有人")
	} else if len(req.AtUsers) > 0 {
		sb.WriteString("\n")
		for _, uid := range req.AtUsers {
			sb.WriteString("@")
			sb.WriteString(uid)
			sb.WriteString(" ")
		}
	}

	return sb.String(), nil
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
