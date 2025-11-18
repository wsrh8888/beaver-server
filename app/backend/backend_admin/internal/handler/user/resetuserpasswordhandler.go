package handler

import (
	"beaver/app/backend/backend_admin/internal/logic/user"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ResetUserPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResetUserPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewResetUserPasswordLogic(r.Context(), svcCtx)
		resp, err := l.ResetUserPassword(&req)
		response.Response(r, w, resp, err)
	}
}
