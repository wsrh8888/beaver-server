package handler

import (
	"beaver/app/open/open_portal/internal/logic/digital"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetDigitalConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDigitalConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetDigitalConfigLogic(r.Context(), svcCtx)
		resp, err := l.GetDigitalConfig(&req)
		response.Response(r, w, resp, err)
	}
}
