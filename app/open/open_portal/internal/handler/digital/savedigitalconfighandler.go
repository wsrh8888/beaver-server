package handler

import (
	"beaver/app/open/open_portal/internal/logic/digital"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SaveDigitalConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveDigitalConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewSaveDigitalConfigLogic(r.Context(), svcCtx)
		resp, err := l.SaveDigitalConfig(&req)
		response.Response(r, w, resp, err)
	}
}
