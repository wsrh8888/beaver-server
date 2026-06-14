package handler

import (
	systemLogic "beaver/app/backend/backend_admin/internal/logic/system"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteAuthorityHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteAuthorityReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := systemLogic.NewDeleteAuthorityLogic(r.Context(), svcCtx)
		resp, err := l.DeleteAuthority(&req)
		response.Response(r, w, resp, err)
	}
}
