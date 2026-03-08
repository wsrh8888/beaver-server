package handler

import (
	"beaver/app/chat/chat_api/internal/logic"
	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func recallMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecallMessageReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		if req.MessageID == "" {
			response.Response(r, w, nil, errors.New("messageId不能为空"))
			return
		}

		l := logic.NewRecallMessageLogic(r.Context(), svcCtx)
		resp, err := l.RecallMessage(&req)
		response.Response(r, w, resp, err)
	}
}
