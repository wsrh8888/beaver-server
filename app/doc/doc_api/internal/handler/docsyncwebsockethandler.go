package handler

import (
	"net/http"

	"beaver/app/doc/doc_api/internal/logic"
	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DocSyncWebsocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WsReq
		if err := httpx.Parse(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		l := logic.NewDocSyncWebsocketLogic(r.Context(), svcCtx)
		_, _ = l.DocSyncWebsocket(&req, w, r)
	}
}
