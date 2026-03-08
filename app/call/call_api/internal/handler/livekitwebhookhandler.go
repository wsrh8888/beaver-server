package handler

import (
	"beaver/app/call/call_api/internal/logic"
	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/common/response"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LiveKitWebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LiveKitWebhookReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 读取 Raw Body 用于签名校验
		body, _ := io.ReadAll(r.Body)
		req.Body = body

		l := logic.NewLiveKitWebhookLogic(r.Context(), svcCtx)
		resp, err := l.LiveKitWebhook(&req)
		response.Response(r, w, resp, err)
	}
}
