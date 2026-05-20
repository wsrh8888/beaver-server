// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"beaver/app/open/open_api/internal/logic/webhook"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 群自定义机器人 Webhook（响应为 code/msg，不走统一 result 包装）
func IncomingWebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IncomingWebhookReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJsonCtx(r.Context(), w, &types.IncomingWebhookRes{
				Code: 403,
				Msg:  "请求参数解析失败",
			})
			return
		}

		l := webhook.NewIncomingWebhookLogic(r.Context(), svcCtx)
		resp, err := l.IncomingWebhook(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
