package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/group"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func MuteGroupMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MuteGroupMemberReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewMuteGroupMemberLogic(r.Context(), svcCtx)
		resp, err := l.MuteGroupMember(&req)
		response.Response(r, w, resp, err)
	}
}
