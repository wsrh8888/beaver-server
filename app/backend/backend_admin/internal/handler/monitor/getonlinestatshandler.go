package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/monitor"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetOnlineStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOnlineStatsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetOnlineStatsLogic(r.Context(), svcCtx)
		resp, err := l.GetOnlineStats(&req)
		response.Response(r, w, resp, err)
	}
}
