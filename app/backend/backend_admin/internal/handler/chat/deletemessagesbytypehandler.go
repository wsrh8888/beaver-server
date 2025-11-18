package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/chat"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteMessagesByTypeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteMessagesByTypeReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewDeleteMessagesByTypeLogic(r.Context(), svcCtx)
		resp, err := l.DeleteMessagesByType(&req)
		response.Response(r, w, resp, err)
	}
}
