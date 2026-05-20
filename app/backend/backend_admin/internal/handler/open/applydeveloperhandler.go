package handler

import (
	"net/http"

	openLogic "beaver/app/backend/backend_admin/internal/logic/open"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ApplyDeveloperHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApplyDeveloperReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := openLogic.NewApplyDeveloperLogic(r.Context(), svcCtx)
		resp, err := l.ApplyDeveloper(&req)
		response.Response(r, w, resp, err)
	}
}
