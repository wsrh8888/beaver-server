package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/friend"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RestoreFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RestoreFriendReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewRestoreFriendLogic(r.Context(), svcCtx)
		resp, err := l.RestoreFriend(&req)
		response.Response(r, w, resp, err)
	}
}
