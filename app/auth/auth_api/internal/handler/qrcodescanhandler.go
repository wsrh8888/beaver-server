package handler

import (
	"errors"
	"net/http"

	"beaver/app/auth/auth_api/internal/logic"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func qrcodeScanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QrcodeScanReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		if req.Token == "" {
			response.Response(r, w, nil, errors.New("二维码 token 不能为空"))
			return
		}
		if req.AuthToken == "" {
			response.Response(r, w, nil, errors.New("未携带登录凭证"))
			return
		}

		l := logic.NewQrcodeScanLogic(r.Context(), svcCtx)
		resp, err := l.QrcodeScan(&req)
		response.Response(r, w, resp, err)
	}
}
