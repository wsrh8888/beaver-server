package handler

import (
	logic "beaver/app/open/open_api/internal/logic/message"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SendRichTextMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendRichTextMessageReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewSendRichTextMessageLogic(r.Context(), svcCtx)
		resp, err := l.SendRichTextMessage(&req)
		response.Response(r, w, resp, err)
	}
}
