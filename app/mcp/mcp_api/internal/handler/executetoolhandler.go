package handler

import (
	"net/http"

	"beaver/app/mcp/mcp_api/internal/logic"
	"beaver/app/mcp/mcp_api/internal/svc"
	"beaver/app/mcp/mcp_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 执行MCP工具
func ExecuteToolHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExecuteToolReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewExecuteToolLogic(r.Context(), svcCtx)
		resp, err := l.ExecuteTool(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
