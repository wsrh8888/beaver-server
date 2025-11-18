package handler

import (
	"beaver/app/group/group_api/internal/logic"
	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func groupJoinRequestSyncHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupJoinRequestSyncReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGroupJoinRequestSyncLogic(r.Context(), svcCtx)
		resp, err := l.GroupJoinRequestSync(&req)
		response.Response(r, w, resp, err)
	}
}
