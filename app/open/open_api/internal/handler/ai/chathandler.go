package handler

import (
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ChatHandler AI 对话流式接口（对标 OpenAI API）
func ChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 设置 SSE 响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no") // Nginx 禁用缓冲

		// 2. 解析请求
		var req types.ChatReq
		if err := httpx.Parse(r, &req); err != nil {
			sendSSEError(w, "参数错误: "+err.Error())
			return
		}

		// 3. 从 context 获取 app_id（中间件注入）
		appID := r.Context().Value("app_id")
		if appID == nil {
			sendSSEError(w, "未认证")
			return
		}

		// 4. 查询 Bot 配置
		l := NewChatLogic(r.Context(), svcCtx)
		err := l.StreamChat(appID.(string), &req, w)
		if err != nil {
			sendSSEError(w, err.Error())
		}
	}
}

// 发送 SSE 错误
func sendSSEError(w http.ResponseWriter, errMsg string) {
	data := map[string]string{
		"error": errMsg,
	}
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(jsonData))
	w.(http.Flusher).Flush()
}
