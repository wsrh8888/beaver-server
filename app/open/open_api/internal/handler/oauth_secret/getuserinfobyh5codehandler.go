package handler

import (
	logic "beaver/app/open/open_api/internal/logic/oauth_secret"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserInfoByH5CodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserInfoByH5CodeReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetUserInfoByH5CodeLogic(r.Context(), svcCtx)
		resp, err := l.GetUserInfoByH5Code(&req)
		response.Response(r, w, resp, err)
	}
}
