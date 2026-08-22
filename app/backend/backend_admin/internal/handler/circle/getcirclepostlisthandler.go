package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/circle"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCirclePostListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCirclePostListReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetCirclePostListLogic(r.Context(), svcCtx)
		resp, err := l.GetCirclePostList(&req)
		response.Response(r, w, resp, err)
	}
}
