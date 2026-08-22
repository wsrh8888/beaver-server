package handler

import (
	logic "beaver/app/platform/platform_api/internal/logic/workbench"
	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListWorkbenchAppsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListWorkbenchAppsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewListWorkbenchAppsLogic(r.Context(), svcCtx)
		resp, err := l.ListWorkbenchApps(&req)
		response.Response(r, w, resp, err)
	}
}
