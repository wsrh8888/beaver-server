package handler

import (
	"beaver/app/doc/doc_api/internal/logic"
	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListCommentsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListCommentsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewListCommentsLogic(r.Context(), svcCtx)
		resp, err := l.ListComments(&req)
		response.Response(r, w, resp, err)
	}
}
