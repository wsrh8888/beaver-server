package handler

import (
	logic "beaver/app/circle/circle_api/internal/logic/circle"
	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ResolveCircleInviteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResolveCircleInviteReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewResolveCircleInviteLogic(r.Context(), svcCtx)
		resp, err := l.ResolveCircleInvite(&req)
		response.Response(r, w, resp, err)
	}
}
