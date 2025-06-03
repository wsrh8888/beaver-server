package core

import (
	"beaver/app/gateway/gateway_api/types"
	"beaver/common/etcd"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	// 执行鉴权接口
	authAddr := etcd.GetServiceAddr(p.Config.Etcd, "auth_api")
	if authAddr == "" {
		return false
	}

	authUrl := fmt.Sprintf("http://%s/api/auth/authentication", authAddr)
	authReq, _ := http.NewRequest("GET", authUrl, nil)
	authReq.Header.Set("ValidPath", req.URL.Path)

	// 添加 User-Agent 头信息
	authReq.Header.Set("User-Agent", req.Header.Get("User-Agent"))

	token := getToken(req)
	if token != "" {
		authReq.Header.Set("Token", token)
	}

	// 设置请求超时
	authClient := &http.Client{Timeout: 10 * time.Second}
	authRes, err := authClient.Do(authReq)
	if err != nil {
		logx.Errorf("认证服务错误: %s", err)
		return false
	}
	defer authRes.Body.Close()

	// 解析响应
	type Response struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Result *struct {
			UserID string `json:"userId"`
		} `json:"result"`
	}

	byteData, err := io.ReadAll(authRes.Body)
	if err != nil {
		logx.Errorf("读取认证服务响应错误: %s", err)
		return false
	}

	var authResponse Response
	authErr := json.Unmarshal(byteData, &authResponse)
	if authErr != nil {
		logx.Errorf("解析认证服务响应错误: %s", authErr)
		return false
	}
	// 检查响应代码
	if authResponse.Code != 0 {
		logx.Errorf("认证服务返回异常: %v", authResponse)
		return false
	}

	if authResponse.Result != nil {
		req.Header.Set("Beaver-User-Id", authResponse.Result.UserID)
	}
	logx.Infof("认证成功: %v", authResponse)
	return true
}

func (p Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	uuid := getUuid(req)
	logx.Info("requst: ", "唯一标识: ", uuid, req.URL.Path)

	// 匹配路由
	regex, _ := regexp.Compile(`/api/(.*?)/`)

	addrList := regex.FindStringSubmatch(req.URL.Path)
	if len(addrList) != 2 {
		writeErrorResponse(res, "请求不匹配", http.StatusBadRequest, uuid)
		return
	}

	service := addrList[1]
	addr := etcd.GetServiceAddr(p.Config.Etcd, service+"_api")
	if addr == "" {
		logx.Error("未匹配到服务")
		writeErrorResponse(res, "未匹配到服务", http.StatusServiceUnavailable, uuid)
		return
	}
	token := getToken(req)
	req.Header.Set("Token", token)
	if !p.auth(res, req) {
		writeErrorResponse(res, "鉴权失败", http.StatusServiceUnavailable, uuid)
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

	return uuid
}
