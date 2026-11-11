package handler

import (
	"beaver/app/doc/doc_api/internal/logic"
	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SetFavoriteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetFavoriteReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewSetFavoriteLogic(r.Context(), svcCtx)
		resp, err := l.SetFavorite(&req)
		response.Response(r, w, resp, err)
	}
}
