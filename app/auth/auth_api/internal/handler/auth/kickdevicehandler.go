package handler

import (
	"beaver/app/auth/auth_api/internal/logic/auth"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func KickDeviceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.KickDeviceReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewKickDeviceLogic(r.Context(), svcCtx)
		resp, err := l.KickDevice(&req)
		response.Response(r, w, resp, err)
	}
}
