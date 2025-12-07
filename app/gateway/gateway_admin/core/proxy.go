package core

import (
	"beaver/app/gateway/gateway_admin/types"
	"beaver/common/etcd"
	"beaver/core"
	"beaver/utils"
	"beaver/utils/jwts"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
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
	Config       types.Config
	backendURL   *url.URL // 解析后的 backend_admin URL
	reverseProxy *httputil.ReverseProxy
	mu           sync.RWMutex  // 保护 backendURL 的并发访问
	lastUpdate   time.Time     // 上次更新后端地址的时间
	redisClient  *redis.Client // Redis 客户端（可选，用于验证 token）
}

func (p *Proxy) ensureRedis() {
	if p.redisClient != nil || p.Config.Redis.Addr == "" {
		return
	}
	p.redisClient = core.InitRedis(p.Config.Redis.Addr, p.Config.Redis.Password, p.Config.Redis.Db)
	logx.Info("Redis 客户端初始化成功，将验证 token 在 Redis 中的有效性")
}

// isWhiteList 检查路径是否在白名单中
func (p *Proxy) isWhiteList(path string) bool {
	for _, whitePath := range p.Config.WhiteList {
		if strings.HasPrefix(path, whitePath) {
			return true
		}
	}
	return false
}

// getBackendAddr 从 etcd 获取后端服务地址
func (p *Proxy) getBackendAddr() (string, error) {
	if p.Config.Etcd == "" {
		return "", fmt.Errorf("未配置 Etcd，无法获取后端服务地址")
	}

	// 固定使用 backend_admin 作为服务 key
	backendKey := "backend_admin"
	addr := etcd.GetServiceAddr(p.Config.Etcd, backendKey)
	if addr == "" {
		return "", fmt.Errorf("从 etcd 获取后端服务地址失败: %s", backendKey)
	}

	return addr, nil
}

// updateBackendAddr 更新后端服务地址（从 etcd 或配置中获取）
func (p *Proxy) updateBackendAddr() error {
	addr, err := p.getBackendAddr()
	if err != nil {
		return err
	}

	backendURL := fmt.Sprintf("http://%s", addr)
	parsedURL, err := url.Parse(backendURL)
	if err != nil {
		return fmt.Errorf("解析后端地址失败: %v", err)
	}

	p.mu.Lock()
	p.backendURL = parsedURL
	p.reverseProxy = httputil.NewSingleHostReverseProxy(parsedURL)
	p.lastUpdate = time.Now()
	p.mu.Unlock()

	// 设置超时
	timeout := 30 * time.Second
	if p.Config.Timeout.Backend > 0 {
		timeout = time.Duration(p.Config.Timeout.Backend) * time.Second
	}
	p.reverseProxy.Transport = &http.Transport{
		ResponseHeaderTimeout: timeout,
	}

	// 修改 Director 函数以保留请求头
	originalDirector := p.reverseProxy.Director
	p.reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		// 保留重要请求头
		if userAgent := req.Header.Get("User-Agent"); userAgent != "" {
			req.Header.Set("User-Agent", userAgent)
		}
		// 设置 X-Forwarded-For
		if clientIP := req.Header.Get("X-Forwarded-For"); clientIP == "" {
			req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		}
		// 设置 X-Real-IP
		req.Header.Set("X-Real-IP", req.RemoteAddr)
		// 保留 Beaver-User-Id（网关认证后设置的用户ID）
		if userId := req.Header.Get("Beaver-User-Id"); userId != "" {
			req.Header.Set("Beaver-User-Id", userId)
		}
	}

	logx.Infof("更新后端服务地址: %s", addr)
	return nil
}

// auth 认证用户（网关直接解析 JWT，避免调用后端认证接口，提升性能）
// 参考大厂和开源 IM 的做法：网关统一鉴权，后端信任网关传递的用户信息
func (p *Proxy) auth(res http.ResponseWriter, req *http.Request) (ok bool) {
	// 1. 检查白名单
	if utils.InListByRegex(p.Config.WhiteList, req.URL.Path) {
		logx.Infof("白名单请求：%s", req.URL.Path)
		return true
	}

	// 2. 获取 token
	token := getToken(req)
	if token == "" {
		logx.Error("token为空")
		return false
	}

	// 3. 直接解析 JWT（避免 HTTP 调用，提升性能）
	if p.Config.Auth.AccessSecret == "" {
		logx.Error("未配置 AccessSecret，无法解析 JWT")
		return false
	}

	claims, err := jwts.ParseToken(token, p.Config.Auth.AccessSecret)
	if err != nil {
		logx.Errorf("JWT解析失败: %v", err)
		return false
	}

	// 4. 可选：验证 token 在 Redis 中的有效性（更安全，防止 token 被撤销）
	if p.redisClient != nil {
		key := fmt.Sprintf("admin_login_%s", claims.UserID)
		storedToken, err := p.redisClient.Get(key).Result()
		if err != nil {
			logx.Errorf("从 Redis 获取 token 失败: %v, userId=%s", err, claims.UserID)
			return false
		}
		if storedToken != token {
			logx.Errorf("token 不一致，可能已被撤销: userId=%s", claims.UserID)
			return false
		}
	}

	// 5. 设置用户ID到请求头（后端服务会信任此 header，不再重复认证）
	req.Header.Set("Beaver-User-Id", claims.UserID)
	logx.Infof("网关JWT验证成功: userId=%s", claims.UserID)

	return true
}

// ServeHTTP 处理 HTTP 请求
func (p *Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	uuid := getUuid(req)
	startTime := time.Now()
	path := req.URL.Path

	logx.Infof("[%s] %s %s", uuid, req.Method, path)

	// 健康检查接口
	if path == "/health" || path == "/ping" {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("ok"))
		return
	}

	// 只处理 /admin/* 路径
	if !strings.HasPrefix(path, "/admin/") {
		writeErrorResponse(res, "请求路径必须以 /admin/ 开头", http.StatusBadRequest, uuid)
		return
	}

	// 检查后端服务是否已初始化，如果未初始化或需要更新，则尝试更新
	if err := p.updateBackendAddr(); err != nil {
		logx.Errorf("[%s] 更新后端服务地址失败: %v", uuid, err)
		writeErrorResponse(res, "网关服务未就绪，请稍后重试", http.StatusServiceUnavailable, uuid)
		return
	}

	// 检查白名单（登录等接口不需要认证）
	needAuth := !p.isWhiteList(path)

	// 设置 Token 到请求头
	token := getToken(req)
	if token != "" {
		req.Header.Set("Token", token)
	}

	// 执行认证（白名单接口跳过）
	if needAuth {
		p.ensureRedis()
		if !p.auth(res, req) {
			writeErrorResponse(res, "认证失败，请先登录", http.StatusUnauthorized, uuid)
			return
		}
	}

	// 记录响应时间
	p.mu.RLock()
	reverseProxy := p.reverseProxy
	p.mu.RUnlock()

	if reverseProxy == nil {
		logx.Errorf("[%s] 反向代理未初始化", uuid)
		writeErrorResponse(res, "网关服务未就绪，请稍后重试", http.StatusServiceUnavailable, uuid)
		return
	}

	originalModifyResponse := reverseProxy.ModifyResponse
	reverseProxy.ModifyResponse = func(resp *http.Response) error {
		duration := time.Since(startTime)
		logx.Infof("[%s] %s %s -> %d (耗时: %v)", uuid, req.Method, path, resp.StatusCode, duration)
		if originalModifyResponse != nil {
			return originalModifyResponse(resp)
		}
		return nil
	}

	// 执行反向代理
	reverseProxy.ServeHTTP(res, req)
}

func logResponseBody(uuid string, res *http.Response) {
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logx.Errorf("读取响应体错误: %s", err)
		return
	}
	res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 恢复响应体
	logx.Info("response: ", "唯一标识: ", uuid, string(bodyBytes))
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
	if uuid == "" {
		uuid = req.Header.Get("X-Request-Id")
	}
	return uuid
}
