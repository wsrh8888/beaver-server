package core

import (
	"beaver/app/gateway/gateway_api/types"
	"beaver/common/etcd"
	"beaver/utils"
	"beaver/utils/jwts"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func writeErrorResponse(res http.ResponseWriter, msg string, statusCode int, uuid string) {
	res.WriteHeader(statusCode)
	res.Header().Set("Content-Type", "application/json")
	response := BaseResponse{Code: 1, Msg: msg}
	byteData, _ := json.Marshal(response)
	res.Write(byteData)
	logx.Info("response: ", "唯一标识: ", uuid, string(byteData))

}

type Proxy struct {
	Config types.Config
}

func (p Proxy) auth(res http.ResponseWriter, req *http.Request) (ok bool) {
	// 1. 检查白名单
	if utils.InListByRegex(p.Config.WhiteList, req.URL.Path) {
		logx.Infof("白名单请求：%s", req.URL.Path)
		return true
	}

	// 2. 判断是否是开放平台请求
	if isOAuthRequest(req.URL.Path) {
		return p.oauthAuth(res, req)
	}

	// 3. 普通用户请求，使用JWT验证
	return p.jwtAuth(res, req)
}

// isOAuthRequest 判断是否是开放平台OAuth请求
func isOAuthRequest(path string) bool {
	return len(path) >= 9 && path[:9] == "/api/open"
}

// oauthAuth OAuth认证（开放平台）
func (p Proxy) oauthAuth(res http.ResponseWriter, req *http.Request) (ok bool) {
	// 获取 Access Token
	token := getToken(req)
	if token == "" {
		logx.Error("OAuth token为空")
		return false
	}

	// TODO: 这里应该调用 open_api 验证 token 的有效性
	// 目前先透传 token，由 open_api 自己验证
	req.Header.Set("Authorization", "Bearer "+token)
	
	logx.Info("OAuth请求，透传token到open_api")
	return true
}

// jwtAuth JWT认证（普通用户）
func (p Proxy) jwtAuth(res http.ResponseWriter, req *http.Request) (ok bool) {
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
	logx.Info("request: ", "唯一标识: ", uuid, req.URL.Path)

	// 限流检查
	if p.Config.Limit.Enable {
		clientIP := getClientIP(req)
		if !p.rateLimitCheck(clientIP) {
			writeErrorResponse(res, "请求频率过高", http.StatusTooManyRequests, uuid)
			return
		}
	}

	token := getToken(req)
	req.Header.Set("Token", token)
	if !p.auth(res, req) {
		writeErrorResponse(res, "网关鉴权失败", http.StatusServiceUnavailable, uuid)
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
		logx.Errorf("未匹配到服务: %s_api", service)
		writeErrorResponse(res, "服务暂时不可用", http.StatusServiceUnavailable, uuid)
		return
	}

	logx.Infof("路由到服务: %s -> %s", service, addr)

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
