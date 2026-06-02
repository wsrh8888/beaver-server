package handler

import (
	logic "beaver/app/auth/auth_api/internal/logic/auth"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func QrcodeScanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QrcodeScanReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewQrcodeScanLogic(r.Context(), svcCtx)
		resp, err := l.QrcodeScan(&req)
		response.Response(r, w, resp, err)
	}
}
