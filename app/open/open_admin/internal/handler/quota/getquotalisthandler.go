package handler

import (
	"beaver/app/open/open_admin/internal/logic/quota"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetQuotaListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetQuotaListReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := quota.NewGetQuotaListLogic(r.Context(), svcCtx)
		resp, err := l.GetQuotaList(&req)
		response.Response(r, w, resp, err)
	}
}
