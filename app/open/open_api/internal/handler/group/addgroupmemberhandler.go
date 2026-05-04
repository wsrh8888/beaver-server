package handler

import (
	logic "beaver/app/open/open_api/internal/logic/group"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddGroupMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddGroupMemberReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewAddGroupMemberLogic(r.Context(), svcCtx)
		resp, err := l.AddGroupMember(&req)
		response.Response(r, w, resp, err)
	}
}
