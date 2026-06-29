package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/workbench"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteWorkbenchAppHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteWorkbenchAppReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewDeleteWorkbenchAppLogic(r.Context(), svcCtx)
		resp, err := l.DeleteWorkbenchApp(&req)
		response.Response(r, w, resp, err)
	}
}
