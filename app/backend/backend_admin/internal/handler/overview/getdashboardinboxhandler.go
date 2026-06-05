package handler

import (
	overviewLogic "beaver/app/backend/backend_admin/internal/logic/overview"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetDashboardInboxHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDashboardInboxReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := overviewLogic.NewGetDashboardInboxLogic(r.Context(), svcCtx)
		resp, err := l.GetDashboardInbox(&req)
		response.Response(r, w, resp, err)
	}
}
