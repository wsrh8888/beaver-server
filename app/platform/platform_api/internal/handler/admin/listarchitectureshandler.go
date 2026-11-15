package handler

import (
	logic "beaver/app/platform/platform_api/internal/logic/admin"
	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListArchitecturesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListArchitecturesReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewListArchitecturesLogic(r.Context(), svcCtx)
		resp, err := l.ListArchitectures(&req)
		response.Response(r, w, resp, err)
	}
}
