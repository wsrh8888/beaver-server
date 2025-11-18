package handler

import (
	"beaver/app/backend/backend_admin/internal/logic/update"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetArchitecturesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetArchitecturesReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetArchitecturesLogic(r.Context(), svcCtx)
		resp, err := l.GetArchitectures(&req)
		response.Response(r, w, resp, err)
	}
}
