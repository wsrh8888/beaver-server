package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/auth"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AuthenticationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthenticationReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if req.Token == "" {
			response.Response(r, w, nil, errors.New("token不能为空"))
			return
		}

		l := logic.NewAuthenticationLogic(r.Context(), svcCtx)
		resp, err := l.Authentication(&req)
		response.Response(r, w, resp, err)
	}
}
