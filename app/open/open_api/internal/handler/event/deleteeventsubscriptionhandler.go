package handler

import (
	logic "beaver/app/open/open_api/internal/logic/event"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteEventSubscriptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteEventSubscriptionReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewDeleteEventSubscriptionLogic(r.Context(), svcCtx)
		resp, err := l.DeleteEventSubscription(&req)
		response.Response(r, w, resp, err)
	}
}
