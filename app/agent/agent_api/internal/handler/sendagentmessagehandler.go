package handler

import (
	"net/http"

	"beaver/app/agent/agent_api/internal/logic"
	"beaver/app/agent/agent_api/internal/svc"
	"beaver/app/agent/agent_api/internal/types"
	"beaver/common/response"
	"beaver/common/sse"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func sendAgentMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendAgentMessageReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		emit, err := sse.Upgrade(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		l := logic.NewSendAgentMessageLogic(r.Context(), svcCtx)
		l.SendAgentMessage(emit, &req)
	}
}
