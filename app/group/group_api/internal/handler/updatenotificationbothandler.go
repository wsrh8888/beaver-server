package handler

import (
	"beaver/app/group/group_api/internal/logic"
	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func updateNotificationBotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateNotificationBotReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewUpdateNotificationBotLogic(r.Context(), svcCtx)
		resp, err := l.UpdateNotificationBot(&req)
		response.Response(r, w, resp, err)
	}
}
