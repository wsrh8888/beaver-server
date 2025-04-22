package handler

import (
	"beaver/app/user/user_api/internal/logic"
	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getUserSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserSettingsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetUserSettingsLogic(r.Context(), svcCtx)
		resp, err := l.GetUserSettings(&req)
		response.Response(r, w, resp, err)
	}
}
