package handler

import (
	"net/http"

	"beaver/app/datasync/datasync_api/internal/logic"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getSyncNotificationReadCursorsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSyncNotificationReadCursorsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetSyncNotificationReadCursorsLogic(r.Context(), svcCtx)
		resp, err := l.GetSyncNotificationReadCursors(&req)
		response.Response(r, w, resp, err)
	}
}

