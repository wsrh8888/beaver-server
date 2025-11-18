package handler

import (
	"beaver/app/backend/backend_admin/internal/logic/track"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetEventListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetEventListReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetEventListLogic(r.Context(), svcCtx)
		resp, err := l.GetEventList(&req)
		response.Response(r, w, resp, err)
	}
}
