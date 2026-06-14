// Package corepush 提供离线 Push（FCM / APNs）。
// WS 在线态见 coreonline；设备档案见 auth_models.AuthDeviceModel。
package corepush

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	PlatformFCM  = "fcm"  // Android / 使用 Firebase SDK 的 iOS
	PlatformAPNs = "apns" // iOS 原生 APNs device token
)

// ─────────────────────────────────────────────────────────────
// 一、对外接口（业务层只需关心这部分）
// ─────────────────────────────────────────────────────────────

// Config 推送配置，由 chat_rpc/svc 从 yaml 组装后传入。
type Config struct {
	Enabled bool
	FCM     FCMConfig  // Android 必填（iOS 若走 Firebase 也只需开 FCM）
	APNs    APNsConfig // iOS 原生 token 时开启，与 FCM 可只开其一
}

// FCMConfig Firebase Cloud Messaging（Google 服务账号 JSON）
type FCMConfig struct {
	Enabled         bool
	ProjectID       string // Firebase 项目 ID
	CredentialsFile string // firebase-adminsdk-xxx.json 路径
}

// APNsConfig Apple Push Notification service（.p8 密钥）
type APNsConfig struct {
	Enabled    bool
	KeyFile    string // AuthKey_XXXX.p8
	KeyID      string
	TeamID     string
	BundleID   string // apns-topic，如 com.beaver.im
	Production bool   // false=沙箱 api.development.push.apple.com
}

// Message 单条推送内容；Data 供客户端唤醒后跳转（conversationId/messageId 等）
type Message struct {
	Title string
	Body  string
	Data  map[string]string
}

// PushSender 离线推送发送器，在 chat_rpc ServiceContext 中初始化一次。
type PushSender struct {
	cfg  Config
	fcm  *fcmClient
	apns *apnsClient
}

// PushToken 单设备 Push 注册信息
type PushToken struct {
	DeviceID     string
	PushToken    string
	PushPlatform string
}

// NewPushSender 创建发送器；凭证读取失败只打日志，不 panic。
func NewPushSender(cfg Config) *PushSender {
	s := &PushSender{cfg: cfg}
	if cfg.Enabled && cfg.FCM.Enabled {
		if c, err := newFCMClient(cfg.FCM.ProjectID, cfg.FCM.CredentialsFile); err != nil {
			logx.Errorf("FCM 初始化失败: %v", err)
		} else {
			s.fcm = c
		}
	}
	if cfg.Enabled && cfg.APNs.Enabled {
		if c, err := newAPNsClient(cfg.APNs); err != nil {
			logx.Errorf("APNs 初始化失败: %v", err)
		} else {
			s.apns = c
		}
	}
	return s
}

// Enabled 是否至少有一个通道（FCM 或 APNs）可用。
func (s *PushSender) Enabled() bool {
	return s != nil && s.cfg.Enabled && (s.fcm != nil || s.apns != nil)
}

// SendToTokens 向指定 Token 列表发通知（chat_rpc 经 AuthRpc 拉取 Token 后调用）。
func (s *PushSender) SendToTokens(ctx context.Context, tokens []PushToken, msg Message) {
	if !s.Enabled() {
		return
	}
	for _, tok := range tokens {
		if err := s.sendToToken(ctx, tok, msg); err != nil {
			logx.Errorf("推送失败: deviceId=%s, err=%v", tok.DeviceID, err)
		}
	}
}

func (s *PushSender) sendToToken(ctx context.Context, tok PushToken, msg Message) error {
	switch tok.PushPlatform {
	case PlatformFCM:
		if s.fcm == nil {
			return nil
		}
		return s.fcm.send(ctx, tok.PushToken, msg)
	case PlatformAPNs:
		if s.apns == nil {
			return nil
		}
		return s.apns.send(ctx, tok.PushToken, msg)
	default:
		return nil
	}
}

// ─────────────────────────────────────────────────────────────
// FCM / APNs 协议实现（内部细节，业务层无需调用）
// 行数多是因为 Google/Apple HTTP API 要求 JWT 签名换 token，无法省略。
// ─────────────────────────────────────────────────────────────

type fcmClient struct {
	projectID   string
	email       string
	privateKey  *rsa.PrivateKey
	tokenURI    string
	httpClient  *http.Client
	mu          sync.Mutex
	accessToken string
	tokenExpiry time.Time
}

func newFCMClient(projectID, credentialsFile string) (*fcmClient, error) {
	if projectID == "" || credentialsFile == "" {
		return nil, errors.New("FCM 配置不完整")
	}
	raw, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}
	var sa struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
		TokenURI    string `json:"token_uri"`
	}
	if err := json.Unmarshal(raw, &sa); err != nil {
		return nil, err
	}
	rsaKey, err := parseRSAPrivateKey(sa.PrivateKey)
	if err != nil {
		return nil, err
	}
	tokenURI := sa.TokenURI
	if tokenURI == "" {
		tokenURI = "https://oauth2.googleapis.com/token"
	}
	return &fcmClient{
		projectID: projectID, email: sa.ClientEmail, privateKey: rsaKey,
		tokenURI: tokenURI, httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// send 调用 FCM HTTP v1：POST .../projects/{id}/messages:send
func (c *fcmClient) send(ctx context.Context, deviceToken string, msg Message) error {
	token, err := c.fetchAccessToken(ctx)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]interface{}{
		"message": map[string]interface{}{
			"token":        deviceToken,
			"notification": map[string]string{"title": msg.Title, "body": msg.Body},
			"data":         msg.Data,
			"android":      map[string]string{"priority": "HIGH"},
		},
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", c.projectID),
		bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("FCM status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return nil
}

// fetchAccessToken 用服务账号私钥签 JWT，向 Google OAuth 换取 access_token（带缓存）
func (c *fcmClient) fetchAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry.Add(-time.Minute)) {
		return c.accessToken, nil
	}
	now := time.Now()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": c.email, "sub": c.email, "aud": c.tokenURI,
		"iat": now.Unix(), "exp": now.Add(time.Hour).Unix(),
		"scope": "https://www.googleapis.com/auth/firebase.messaging",
	})
	signed, err := jwtToken.SignedString(c.privateKey)
	if err != nil {
		return "", err
	}
	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Set("assertion", signed)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenURI, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	c.accessToken = result.AccessToken
	c.tokenExpiry = now.Add(time.Duration(result.ExpiresIn) * time.Second)
	return c.accessToken, nil
}

type apnsClient struct {
	keyID, teamID, bundleID, host string
	privateKey                    *ecdsa.PrivateKey
	httpClient                    *http.Client
	mu                            sync.Mutex
	authToken                     string
	tokenExp                      time.Time
}

func newAPNsClient(cfg APNsConfig) (*apnsClient, error) {
	if cfg.KeyFile == "" || cfg.KeyID == "" || cfg.TeamID == "" || cfg.BundleID == "" {
		return nil, errors.New("APNs 配置不完整")
	}
	raw, err := os.ReadFile(cfg.KeyFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(raw)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ecKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("APNs 密钥类型错误")
	}
	host := "https://api.development.push.apple.com"
	if cfg.Production {
		host = "https://api.push.apple.com"
	}
	return &apnsClient{
		keyID: cfg.KeyID, teamID: cfg.TeamID, bundleID: cfg.BundleID, host: host,
		privateKey: ecKey, httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// send 调用 APNs HTTP/2 API：POST /3/device/{deviceToken}
func (c *apnsClient) send(ctx context.Context, deviceToken string, msg Message) error {
	auth, err := c.fetchAuthToken()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert": map[string]string{"title": msg.Title, "body": msg.Body},
			"sound": "default",
		},
	}
	for k, v := range msg.Data {
		payload[k] = v
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.host+"/3/device/"+deviceToken, bytes.NewReader(body))
	req.Header.Set("authorization", "bearer "+auth)
	req.Header.Set("apns-topic", c.bundleID)
	req.Header.Set("apns-push-type", "alert")
	req.Header.Set("apns-priority", "10")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("APNs status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return nil
}

// fetchAuthToken 用 .p8 密钥签 ES256 JWT 作为 APNs 鉴权（带缓存，约 50 分钟有效）
func (c *apnsClient) fetchAuthToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.authToken != "" && time.Now().Before(c.tokenExp.Add(-5*time.Minute)) {
		return c.authToken, nil
	}
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{"iss": c.teamID, "iat": now.Unix()})
	token.Header["kid"] = c.keyID
	signed, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", err
	}
	c.authToken = signed
	c.tokenExp = now.Add(50 * time.Minute)
	return c.authToken, nil
}

func parseRSAPrivateKey(pemText string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemText))
	if block == nil {
		return nil, errors.New("私钥格式无效")
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("非 RSA 私钥")
		}
		return rsaKey, nil
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
