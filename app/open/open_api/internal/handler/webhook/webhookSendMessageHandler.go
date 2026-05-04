// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package webhook

import (
	"net/http"

	"beaver/app/open/open_api/internal/logic/webhook"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 通过 Webhook 发送消息（无需鉴权，通过 URL 中的 token 验证）
func WebhookSendMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WebhookSendMessageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := webhook.NewWebhookSendMessageLogic(r.Context(), svcCtx)
		resp, err := l.WebhookSendMessage(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
