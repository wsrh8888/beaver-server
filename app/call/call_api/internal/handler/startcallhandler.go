package handler

import (
	"beaver/app/call/call_api/internal/logic"
	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func StartCallHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StartCallReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewStartCallLogic(r.Context(), svcCtx)
		resp, err := l.StartCall(&req)
		response.Response(r, w, resp, err)
	}
}
