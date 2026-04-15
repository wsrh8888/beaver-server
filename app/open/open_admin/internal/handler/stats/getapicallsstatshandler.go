package handler

import (
	"beaver/app/open/open_admin/internal/logic/stats"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAPICallsStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAPICallsStatsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := stats.NewGetAPICallsStatsLogic(r.Context(), svcCtx)
		resp, err := l.GetAPICallsStats(&req)
		response.Response(r, w, resp, err)
	}
}
