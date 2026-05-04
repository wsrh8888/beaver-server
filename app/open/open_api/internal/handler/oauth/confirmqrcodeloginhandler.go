package oauth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"beaver/app/open/open_api/internal/logic/oauth"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
)

// 确认扫码登录
func ConfirmQrCodeLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConfirmQrCodeLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := oauth.NewConfirmQrCodeLoginLogic(r.Context(), svcCtx)
		resp, err := l.ConfirmQrCodeLogin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
