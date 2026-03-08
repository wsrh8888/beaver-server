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

func forwardMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ForwardMessageReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		if len(req.MessageIDs) == 0 {
			response.Response(r, w, nil, errors.New("messageIds不能为空"))
			return
		}
		if req.TargetID == "" {
			response.Response(r, w, nil, errors.New("targetId不能为空"))
			return
		}
		if req.ForwardMode != 1 && req.ForwardMode != 2 {
			response.Response(r, w, nil, errors.New("无效的转发模式"))
			return
		}

		l := logic.NewForwardMessageLogic(r.Context(), svcCtx)
		resp, err := l.ForwardMessage(&req)
		response.Response(r, w, resp, err)
	}
}
