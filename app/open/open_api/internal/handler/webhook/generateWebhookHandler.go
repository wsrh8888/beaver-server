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

// 生成群机器人 Webhook URL（对标钉钉/企业微信）
func GenerateWebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenerateWebhookReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := webhook.NewGenerateWebhookLogic(r.Context(), svcCtx)
		resp, err := l.GenerateWebhook(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
