package handler

import (
	"beaver/app/datasync/datasync_api/internal/logic"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getSyncGroupMembersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSyncGroupMembersReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetSyncGroupMembersLogic(r.Context(), svcCtx)
		resp, err := l.GetSyncGroupMembers(&req)
		response.Response(r, w, resp, err)
	}
}
