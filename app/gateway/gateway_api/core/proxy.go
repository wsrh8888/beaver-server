package core

import (
	"beaver/app/gateway/gateway_api/types"
	"beaver/common/etcd"
	"beaver/utils/jwts"
	utils "beaver/utils/list"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var gatewayLog = logger.New("proxy")

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func writeErrorResponse(res http.ResponseWriter, msg string, statusCode int, uuid string) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	response := BaseResponse{Code: 1, Msg: msg}
	byteData, _ := json.Marshal(response)
	res.Write(byteData)
	logx.Info("response: ", "唯一标识: ", uuid, string(byteData))

}

type Proxy struct {
	Config types.Config
}

func (p Proxy) auth(req *http.Request) (ok bool, errMsg string) {
	path := req.URL.Path

	// 1. 公开接口，无需鉴权
	if utils.InListByRegex(p.Config.PublicList, path) {
		return true, ""
	}

	// 2. 自定义鉴权，透传到下游服务 middleware 处理
	if utils.InListByRegex(p.Config.CustomAuthList, path) {
		return true, ""
	}

	// 3. open oauth_secret：Gateway 校验 App-Id / App-Secret 请求头
	if strings.HasPrefix(path, "/api/open/oauth_secret/") {
		return p.oauthSecretAuth(req)
	}

	// 4. open_api：Gateway 默认不鉴权，各接口 logic 自行校验；仅 Beaver JWT 路由例外
	if isOpenApiPassThrough(path) {
		return true, ""
	}

	// 5. 统一 JWT 鉴权
	if !p.jwtAuth(req) {
		return false, "网关鉴权失败"
	}
	return true, ""
}

func (p Proxy) oauthSecretAuth(req *http.Request) (bool, string) {
	appID := req.Header.Get("App-Id")
	if appID == "" {
		return false, "缺少 App-Id 请求头"
	}

	appSecret := req.Header.Get("App-Secret")
	if appSecret == "" {
		return false, "缺少 App-Secret 请求头"
	}
	return true, ""
}

var openApiJwtRoutes = []string{
	`/api/open/oauth/v1/h5_authcode`,
	`/api/open/oauth/v1/qrcode_scan`,
	`/api/open/oauth/v1/qrcode_confirm`,
	`/api/open/oauth/v1/qrcode_cancel`,
}

func isOpenApiPassThrough(path string) bool {
	if !strings.HasPrefix(path, "/api/open/") {
		return false
	}
	if utils.InListByRegex(openApiJwtRoutes, path) {
		return false
	}
	return true
}

// jwtAuth JWT认证（普通用户）
func (p Proxy) jwtAuth(req *http.Request) bool {
	// 获取token
	token := getToken(req)
	if token == "" {
		logx.Error("token为空")
		return false
	}

	// 直接解析JWT（避免HTTP调用）
	claims, err := jwts.ParseToken(token, p.Config.Auth.AccessSecret)
	if err != nil {
		logx.Errorf("JWT解析失败: %v", err)
		return false
	}

	// 设置用户ID和设备ID到请求头
	req.Header.Set("Beaver-User-Id", claims.UserID)

	// 从请求头获取设备ID
	deviceId := req.Header.Get("deviceId")
	if deviceId != "" {
		req.Header.Set("Beaver-Device-Id", deviceId)
	}
	version := req.Header.Get("version")
	if version != "" {
		req.Header.Set("Version", version)
	}

	logx.Infof("JWT验证成功: 用户=%s, 设备=%s", claims.UserID, deviceId)

	return true
}

func (p Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	uuid := getUuid(req)

	// 限流检查
	if p.Config.Limit.Enable {
		clientIP := getClientIP(req)
		if !p.rateLimitCheck(clientIP) {
			gatewayLog.Warn(model.LogMsg{
				Text: "请求频率过高",
				Data: map[string]interface{}{
					"clientIp": clientIP,
					"path":     req.URL.Path,
					"uuid":     uuid,
				},
			})
			writeErrorResponse(res, "请求频率过高", http.StatusTooManyRequests, uuid)
			return
		}
	}

	token := getToken(req)
	req.Header.Set("Token", token)
	ok, authErrMsg := p.auth(req)
	if !ok {
		gatewayLog.Warn(model.LogMsg{
			Text: "网关鉴权失败",
			Data: map[string]interface{}{
				"path": req.URL.Path,
				"uuid": uuid,
				"msg":  authErrMsg,
			},
		})
		writeErrorResponse(res, authErrMsg, http.StatusUnauthorized, uuid)
		return
	}

	// 匹配路由
	regex, _ := regexp.Compile(`/api/(.*?)/`)

	addrList := regex.FindStringSubmatch(req.URL.Path)
	if len(addrList) != 2 {
		writeErrorResponse(res, "请求不匹配", http.StatusBadRequest, uuid)
		return
	}

	service := addrList[1]

	// 增加重试机制
	var addr string
	for i := 0; i < 3; i++ {
		addr = etcd.GetServiceAddr(p.Config.Etcd, service+"_api")
		if addr != "" {
			break
		}
		logx.Errorf("第%d次获取服务地址失败: %s_api", i+1, service)
		time.Sleep(100 * time.Millisecond)
	}

	if addr == "" {
		gatewayLog.Error(model.LogMsg{
			Text: "服务不可用",
			Data: map[string]interface{}{
				"service": service + "_api",
				"path":    req.URL.Path,
				"uuid":    uuid,
			},
		})
		logx.Errorf("未匹配到服务: %s_api", service)
		writeErrorResponse(res, "服务暂时不可用", http.StatusServiceUnavailable, uuid)
		return
	}

	remote, _ := url.Parse(fmt.Sprintf("http://%s", addr))
	reverseProxy := httputil.NewSingleHostReverseProxy(remote)

	// 修改默认的 Director 函数以保留 User-Agent
	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		// 确保 User-Agent 被保留
		if userAgent := req.Header.Get("User-Agent"); userAgent != "" {
			req.Header.Set("User-Agent", userAgent)
		}
	}

	reverseProxy.ServeHTTP(res, req)
}

func getToken(req *http.Request) string {
	token := req.Header.Get("Token")
	if token == "" {
		token = req.URL.Query().Get("token")
	}
	return token
}

func getUuid(req *http.Request) string {
	uuid := req.Header.Get("Uuid")

	return uuid
}

// getClientIP 获取客户端真实IP
func getClientIP(req *http.Request) string {
	// 优先从 X-Forwarded-For 获取
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// 其次从 X-Real-IP 获取
	if xri := req.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// 最后从 RemoteAddr 获取
	return req.RemoteAddr
}

// rateLimitCheck 限流检查
var globalRateLimiter *RateLimiter

func (p Proxy) rateLimitCheck(clientIP string) bool {
	if globalRateLimiter == nil {
		globalRateLimiter = NewRateLimiter(p.Config.Limit.Rate, p.Config.Limit.Burst)
	}
	return globalRateLimiter.Allow(clientIP)
}
