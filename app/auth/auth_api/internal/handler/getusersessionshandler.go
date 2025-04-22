package handler

import (
	"beaver/app/auth/auth_api/internal/logic"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getUserSessionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserSessionsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetUserSessionsLogic(r.Context(), svcCtx)
		resp, err := l.GetUserSessions(&req)
		response.Response(r, w, resp, err)
	}
}
