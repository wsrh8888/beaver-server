package handler

import (
	"beaver/app/notification/notification_api/internal/logic"
	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getUnreadSummaryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUnreadSummaryReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetUnreadSummaryLogic(r.Context(), svcCtx)
		resp, err := l.GetUnreadSummary(&req)
		response.Response(r, w, resp, err)
	}
}
