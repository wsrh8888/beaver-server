package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/group"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetGroupDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGroupDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetGroupDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetGroupDetail(&req)
		response.Response(r, w, resp, err)
	}
}
