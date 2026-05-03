package handler

import (
	"beaver/app/open/open_admin/internal/logic/quota"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConfigQuotaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConfigQuotaReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := quota.NewConfigQuotaLogic(r.Context(), svcCtx)
		resp, err := l.ConfigQuota(&req)
		response.Response(r, w, resp, err)
	}
}
