package handler

import (
	operationsLogic "beaver/app/backend/backend_admin/internal/logic/operations"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserOperationsProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserOperationsProfileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := operationsLogic.NewGetUserOperationsProfileLogic(r.Context(), svcCtx)
		resp, err := l.GetUserOperationsProfile(&req)
		response.Response(r, w, resp, err)
	}
}
