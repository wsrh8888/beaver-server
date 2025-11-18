package handler

import (
	"beaver/app/backend/backend_admin/internal/logic/system"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewDeleteMenuLogic(r.Context(), svcCtx)
		resp, err := l.DeleteMenu(&req)
		response.Response(r, w, resp, err)
	}
}
