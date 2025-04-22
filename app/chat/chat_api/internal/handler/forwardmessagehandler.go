package handler

import (
	"beaver/app/chat/chat_api/internal/logic"
	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func forwardMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ForwardMessageReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewForwardMessageLogic(r.Context(), svcCtx)
		resp, err := l.ForwardMessage(&req)
		response.Response(r, w, resp, err)
	}
}
