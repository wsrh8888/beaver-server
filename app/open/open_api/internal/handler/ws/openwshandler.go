package handler

import (
	logic "beaver/app/open/open_api/internal/logic/ws"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func OpenWsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenWsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewOpenWsLogic(r.Context(), svcCtx)
		resp, err := l.OpenWs(&req, w, r)

		// 连接正常结束时已完成 WebSocket 升级并接管了响应，
		// 此时不能再写 HTTP 响应体，否则会报 hijacked connection。
		if err != nil {
			response.Response(r, w, resp, err)
		}
	}
}
