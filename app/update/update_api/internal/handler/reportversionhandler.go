package handler

import (
	"beaver/app/update/update_api/internal/logic"
	"beaver/app/update/update_api/internal/svc"
	"beaver/app/update/update_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func reportVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReportVersionReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewReportVersionLogic(r.Context(), svcCtx)
		resp, err := l.ReportVersion(&req)
		response.Response(r, w, resp, err)
	}
}
