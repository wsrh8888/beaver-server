package handler

import (
	logic "beaver/app/open/open_api/internal/logic/oauth_secret"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RevokeTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RevokeTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewRevokeTokenLogic(r.Context(), svcCtx)
		resp, err := l.RevokeToken(&req, r.Header.Get("App-Id"), r.Header.Get("App-Secret"))
		response.Response(r, w, resp, err)
	}
}
