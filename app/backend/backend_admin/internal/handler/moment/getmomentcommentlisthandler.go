package handler

import (
	"beaver/app/backend/backend_admin/internal/logic/moment"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMomentCommentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMomentCommentListReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetMomentCommentListLogic(r.Context(), svcCtx)
		resp, err := l.GetMomentCommentList(&req)
		response.Response(r, w, resp, err)
	}
}
