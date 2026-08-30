package handler

import (
	"beaver/app/agent/agent_api/internal/logic"
	"beaver/app/agent/agent_api/internal/svc"
	"beaver/app/agent/agent_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func listAgentModelsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListAgentModelsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewListAgentModelsLogic(r.Context(), svcCtx)
		resp, err := l.ListAgentModels(&req)
		response.Response(r, w, resp, err)
	}
}
