package handler

import (
	"beaver/app/open/open_admin/internal/logic/permission"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAppPermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAppPermissionsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := permission.NewGetAppPermissionsLogic(r.Context(), svcCtx)
		resp, err := l.GetAppPermissions(&req)
		response.Response(r, w, resp, err)
	}
}
