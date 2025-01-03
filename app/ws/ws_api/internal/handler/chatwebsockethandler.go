package handler

import (
	"net/http"

	"beaver/app/ws/ws_api/internal/logic"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func chatWebsocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewChatWebsocketLogic(r.Context(), svcCtx)

		resp, err := l.ChatWebsocket(&req, w, r)

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}