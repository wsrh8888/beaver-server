package handler

import (
	logic "beaver/app/open/open_admin/internal/logic/app"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAppDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAppDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetAppDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetAppDetail(&req)
		response.Response(r, w, resp, err)
	}
}
