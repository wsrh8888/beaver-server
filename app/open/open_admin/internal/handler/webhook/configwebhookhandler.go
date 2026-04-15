package handler

import (
	"beaver/app/open/open_admin/internal/logic/webhook"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConfigWebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConfigWebhookReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := webhook.NewConfigWebhookLogic(r.Context(), svcCtx)
		resp, err := l.ConfigWebhook(&req)
		response.Response(r, w, resp, err)
	}
}
