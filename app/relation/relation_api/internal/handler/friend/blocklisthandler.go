package handler

import (
	"beaver/app/relation/relation_api/internal/logic/friend"
	"beaver/app/relation/relation_api/internal/svc"
	"beaver/app/relation/relation_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func BlockListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BlockListReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewBlockListLogic(r.Context(), svcCtx)
		resp, err := l.BlockList(&req)
		response.Response(r, w, resp, err)
	}
}
