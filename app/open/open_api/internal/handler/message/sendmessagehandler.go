package handler

import (
	"beaver/app/open/open_api/internal/logic/message"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SendMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendMessageReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := message.NewSendMessageLogic(r.Context(), svcCtx)
		resp, err := l.SendMessage(&req)
		response.Response(r, w, resp, err)
	}
}
