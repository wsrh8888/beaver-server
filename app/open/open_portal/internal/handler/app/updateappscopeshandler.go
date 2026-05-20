package handler

import (
	logic "beaver/app/open/open_portal/internal/logic/app"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateAppScopesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateAppScopesReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewUpdateAppScopesLogic(r.Context(), svcCtx)
		resp, err := l.UpdateAppScopes(&req)
		response.Response(r, w, resp, err)
	}
}
