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

func qrcodeStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QrcodeStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		if req.Token == "" {
			response.Response(r, w, nil, errors.New("token 不能为空"))
			return
		}

		l := logic.NewQrcodeStatusLogic(r.Context(), svcCtx)
		resp, err := l.QrcodeStatus(&req)
		response.Response(r, w, resp, err)
	}
}
