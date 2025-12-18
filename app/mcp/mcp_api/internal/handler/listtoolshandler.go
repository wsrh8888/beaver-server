package handler

import (
	"net/http"

	"beaver/app/mcp/mcp_api/internal/logic"
	"beaver/app/mcp/mcp_api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取所有可用工具列表
func ListToolsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewListToolsLogic(r.Context(), svcCtx)
		resp, err := l.ListTools()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
