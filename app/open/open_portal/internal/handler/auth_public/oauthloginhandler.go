package handler

import (
	logic "beaver/app/open/open_portal/internal/logic/auth_public"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func OAuthLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OAuthLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewOAuthLoginLogic(r.Context(), svcCtx)
		resp, err := l.OAuthLogin(&req)
		response.Response(r, w, resp, err)
	}
}
