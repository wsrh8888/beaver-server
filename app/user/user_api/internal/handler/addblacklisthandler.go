package handler

import (
	"beaver/app/user/user_api/internal/logic"
	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func addBlacklistHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddBlacklistReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewAddBlacklistLogic(r.Context(), svcCtx)
		resp, err := l.AddBlacklist(&req)
		response.Response(r, w, resp, err)
	}
}
