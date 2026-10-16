package agent

import (
	"net/http"

	agentLogic "beaver/app/backend/backend_admin/internal/logic/agent"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAiAssistantHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAiAssistantReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		l := agentLogic.NewGetAiAssistantLogic(r.Context(), svcCtx)
		resp, err := l.GetAiAssistant(&req)
		response.Response(r, w, resp, err)
	}
}
