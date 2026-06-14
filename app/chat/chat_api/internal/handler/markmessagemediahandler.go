package handler

import (
	"errors"
	"net/http"

	"beaver/app/chat/chat_api/internal/logic"
	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func markMessageMediaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MarkMessageMediaReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		if len(req.MessageIDs) == 0 {
			response.Response(r, w, nil, errors.New("messageIds不能为空"))
			return
		}
		if len(req.MessageIDs) > 100 {
			response.Response(r, w, nil, errors.New("单次最多批量标记100条消息"))
			return
		}

		l := logic.NewMarkMessageMediaLogic(r.Context(), svcCtx)
		resp, err := l.MarkMessageMedia(&req)
		response.Response(r, w, resp, err)
	}
}
