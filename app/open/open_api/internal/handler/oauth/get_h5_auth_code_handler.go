// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package oauth

import (
	"net/http"

	"beaver/app/open/open_api/internal/logic/oauth"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// H5 免登获取 authCode（需在 WebView 环境中调用）
func GetH5AuthCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetH5AuthCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := oauth.NewGetH5AuthCodeLogic(r.Context(), svcCtx)
		resp, err := l.GetH5AuthCode(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
