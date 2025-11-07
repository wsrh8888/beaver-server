package handler

import (
	"beaver/app/datasync/datasync_api/internal/logic"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getSyncAllUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSyncAllUsersReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetSyncAllUsersLogic(r.Context(), svcCtx)
		resp, err := l.GetSyncAllUsers(&req)
		response.Response(r, w, resp, err)
	}
}
