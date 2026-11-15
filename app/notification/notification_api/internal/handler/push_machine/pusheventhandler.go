package handler

import (
	logic "beaver/app/notification/notification_api/internal/logic/push_machine"
	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PushEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PushEventReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewPushEventLogic(r.Context(), svcCtx)
		resp, err := l.PushEvent(&req)
		response.Response(r, w, resp, err)
	}
}
